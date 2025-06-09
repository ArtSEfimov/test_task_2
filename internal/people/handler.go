package people

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"go_test_task_2/config"
	"go_test_task_2/pkg/response"
	"net/http"
	"strconv"
	"strings"
)

type HandlerDeps struct {
	Config     *config.Config
	Repository *Repository
}

func NewHandlerDeps(conf *config.Config, repository *Repository) *HandlerDeps {
	return &HandlerDeps{
		Config:     conf,
		Repository: repository,
	}
}

type Handler struct {
	config     *config.Config
	repository *Repository
}

func NewHandler(router *http.ServeMux, deps *HandlerDeps) {
	handler := &Handler{
		config:     deps.Config,
		repository: deps.Repository,
	}

	router.HandleFunc("GET /people", handler.GetAll())
	router.HandleFunc("GET /people/{id}", handler.GetByID())
	router.HandleFunc("POST /people", handler.Create())
	router.HandleFunc("PUT  /people/{id}", handler.Update())
	router.HandleFunc("DELETE  /people/{id}", handler.Delete())
}

func (handler *Handler) listLimit(w http.ResponseWriter, sizeParam string) {
	handler.config.InfoLogger.Printf("Received request with sizeParam: %s", sizeParam)

	parseUint, err := strconv.ParseUint(sizeParam, 10, 64)
	if err != nil {
		handler.config.DebugLogger.Printf("Error parsing sizeParam: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := parseUint

	handler.config.InfoLogger.Printf("Constructed query with LIMIT: %d", limit)

	params := "LIMIT $1"
	dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
	query := fmt.Sprintf("SELECT * FROM%s\nORDER BY id\n%s", dbName, params)

	handler.config.InfoLogger.Printf("Using query: %s", query)

	var people AllPeopleResponse
	getErr := handler.repository.Get(query, &people, limit)

	if getErr != nil {
		handler.config.DebugLogger.Printf("Error executing query: %v", getErr)
		http.Error(w, getErr.Error(), http.StatusInternalServerError)
		return
	}

	handler.config.InfoLogger.Printf("Query executed successfully, found %d records", len(people.People))

	response.Json(w, &people, http.StatusOK)

	handler.config.InfoLogger.Println("Finished processing listLimit request")

}

func (handler *Handler) listLimitOffset(w http.ResponseWriter, sizeParam, pageParam string) {

	handler.config.InfoLogger.Printf("Received request with sizeParam: %s, pageParam: %s", sizeParam, pageParam)

	parseUint, err := strconv.ParseUint(pageParam, 10, 64)
	if err != nil {
		handler.config.DebugLogger.Printf("Error parsing pageParam: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	page := parseUint

	parseUint, err = strconv.ParseUint(sizeParam, 10, 64)
	if err != nil {
		handler.config.DebugLogger.Printf("Error parsing sizeParam: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := parseUint

	offset := (page - 1) * limit

	handler.config.InfoLogger.Printf("Calculated offset: %d for page: %d and limit: %d", offset, page, limit)

	params := "LIMIT $1 OFFSET $2"
	dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
	query := fmt.Sprintf("SELECT * FROM%s\nORDER BY id\n%s", dbName, params)

	handler.config.InfoLogger.Printf("Using query: %s", query)

	var people AllPeopleResponse
	getErr := handler.repository.Get(query, &people, limit, offset)

	if getErr != nil {
		handler.config.DebugLogger.Printf("Error executing query: %v", getErr)
		http.Error(w, getErr.Error(), http.StatusInternalServerError)
		return
	}

	handler.config.InfoLogger.Printf("Query executed successfully, found %d records", len(people.People))

	response.Json(w, &people, http.StatusOK)

	handler.config.InfoLogger.Println("Finished processing listLimitOffset request")

}

func (handler *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pageParam := parseQuery(r, "page")
		sizeParam := parseQuery(r, "limit")

		if pageParam != "" && sizeParam != "" {
			handler.config.InfoLogger.Println("Handling request with pagination and limit.")
			handler.listLimitOffset(w, sizeParam, pageParam)
			return
		}

		if sizeParam != "" {
			handler.config.InfoLogger.Println("Handling request with limit only.")
			handler.listLimit(w, sizeParam)
			return
		}

		handler.config.InfoLogger.Printf("Handling request with pageParam: %s, sizeParam: %s", pageParam, sizeParam)

		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("SELECT * FROM%s\nORDER BY id", dbName)

		handler.config.InfoLogger.Printf("Using query: %s", query)

		var people AllPeopleResponse
		dbErr := handler.repository.Get(query, &people)

		if dbErr != nil {
			handler.config.DebugLogger.Printf("Error executing query: %v", dbErr)
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		handler.config.InfoLogger.Printf("Query executed successfully, found %d records", len(people.People))

		response.Json(w, &people, http.StatusOK)

		handler.config.InfoLogger.Println("Finished processing GET request")

	}
}

func (handler *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")

		handler.config.InfoLogger.Printf("Received request for ID: %s", idString)

		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			handler.config.DebugLogger.Printf("Error parsing ID: %v", parseErr)
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
		}

		handler.config.InfoLogger.Printf("Parsed ID successfully: %d", id)

		params := "WHERE id = $1"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("SELECT * FROM%s\n%s", dbName, params)

		handler.config.InfoLogger.Printf("Executing query: %s", query)

		var person Person
		dbErr := handler.repository.GetByID(query, &person, id)

		if dbErr != nil {
			handler.config.DebugLogger.Printf("Error executing query: %v", dbErr)
			var errorNotFound *ErrorNotFound
			if errors.As(dbErr, &errorNotFound) {
				http.Error(w, dbErr.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		handler.config.InfoLogger.Printf("Successfully retrieved person with ID: %d", id)

		response.Json(w, &person, http.StatusOK)

		handler.config.InfoLogger.Println("Finished processing GET request")

	}
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var person Person
		var request Request

		handler.config.InfoLogger.Println("Received request to create a person")

		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&request)
		if decodeErr != nil {
			handler.config.DebugLogger.Printf("Error decoding request body: %v", decodeErr)
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		}
		errorMessages := IsValid(&request)
		if errorMessages != nil {
			handler.config.DebugLogger.Printf("Validation failed for request: %s", strings.Join(errorMessages, ", "))
			http.Error(w, fmt.Sprintf("Validation failed: %s", strings.Join(errorMessages, ", ")), http.StatusBadRequest)
			return
		}

		person.Name = request.Name
		person.Surname = request.Surname
		person.Patronymic = request.Patronymic

		handler.config.InfoLogger.Println("Enriching person data")

		promise := enrichPerson(&person)

		handler.config.InfoLogger.Println("Person data enriched successfully")

		params := "VALUES ($1, $2, $3, $4, $5, $6)"
		returning := "RETURNING id, created_at, updated_at"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("INSERT INTO%s (name, surname, patronymic, age, gender, nationality) %s %s", dbName, params, returning)

		handler.config.InfoLogger.Printf("Executing query: %s", query)

		<-promise

		dbErr := handler.repository.Create(query, &person)
		if dbErr != nil {
			handler.config.DebugLogger.Printf("Error executing query: %v", dbErr)
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		handler.config.InfoLogger.Printf("Successfully created person with ID: %d", person.ID)

		response.Json(w, &person, http.StatusCreated)

		handler.config.InfoLogger.Println("Finished processing POST request")

	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var person Person
		var request Request

		handler.config.InfoLogger.Println("Received request to update a person")

		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			handler.config.DebugLogger.Printf("Error parsing ID: %v", parseErr)
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
		}

		handler.config.InfoLogger.Printf("Parsed ID successfully: %d", id)

		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&request)
		if decodeErr != nil {
			handler.config.DebugLogger.Printf("Error decoding request body: %v", decodeErr)
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		}
		errorMessages := IsValid(&request)
		if errorMessages != nil {
			handler.config.DebugLogger.Printf("Validation failed for request: %s", strings.Join(errorMessages, ", "))
			http.Error(w, fmt.Sprintf("Validation failed: %s", strings.Join(errorMessages, ", ")), http.StatusBadRequest)
			return
		}

		person.Name = request.Name
		person.Surname = request.Surname
		person.Patronymic = request.Patronymic

		handler.config.InfoLogger.Println("Enriching person data")

		promise := enrichPerson(&person)

		handler.config.InfoLogger.Println("Person data enriched successfully")

		params := "SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6"
		selectByID := "WHERE id = $7"
		returning := "RETURNING id, created_at, updated_at"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("UPDATE%s\n%s\n%s\n%s", dbName, params, selectByID, returning)

		handler.config.InfoLogger.Printf("Executing query: %s", query)

		<-promise

		dbErr := handler.repository.Update(query, &person, id)
		if dbErr != nil {
			handler.config.DebugLogger.Printf("Error executing query: %v", dbErr)
			var errorNotFound *ErrorNotFound
			if errors.As(dbErr, &errorNotFound) {
				http.Error(w, dbErr.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		handler.config.InfoLogger.Printf("Successfully updated person with ID: %d", id)

		response.Json(w, &person, http.StatusOK)

		handler.config.InfoLogger.Println("Finished processing PUT request")
	}
}

func (handler *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		handler.config.InfoLogger.Println("Received request to delete a person")

		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			handler.config.DebugLogger.Printf("Error parsing ID: %v", parseErr)
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
		}

		handler.config.InfoLogger.Printf("Parsed ID successfully: %d", id)

		params := "WHERE id = $1"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("DELETE FROM%s\n%s", dbName, params)

		handler.config.InfoLogger.Printf("Executing query: %s", query)

		dbErr := handler.repository.Delete(query, id)

		if dbErr != nil {
			handler.config.DebugLogger.Printf("Error executing query: %v", dbErr)
			var errorNotFound *ErrorNotFound
			if errors.As(dbErr, &errorNotFound) {
				http.Error(w, dbErr.Error(), http.StatusBadRequest)
				return
			}
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		handler.config.InfoLogger.Printf("Successfully deleted person with ID: %d", id)

		response.Json(w, nil, http.StatusNoContent)

		handler.config.InfoLogger.Println("Finished processing DELETE request")

	}
}
