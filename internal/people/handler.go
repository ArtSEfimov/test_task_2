package people

import (
	"bufio"
	"encoding/json"
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

	parseUint, err := strconv.ParseUint(sizeParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := parseUint

	params := "LIMIT $1"
	dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
	query := fmt.Sprintf("SELECT * FROM%s\nORDER BY id\n%s", dbName, params)

	var people AllPeopleResponse
	getErr := handler.repository.Get(query, &people, limit)

	if getErr != nil {
		http.Error(w, getErr.Error(), http.StatusInternalServerError)
		return
	}

	response.Json(w, &people, http.StatusOK)

}

func (handler *Handler) listLimitOffset(w http.ResponseWriter, sizeParam, pageParam string) {

	parseUint, err := strconv.ParseUint(pageParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	page := parseUint

	parseUint, err = strconv.ParseUint(sizeParam, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	limit := parseUint

	offset := (page - 1) * limit

	params := "LIMIT $1 OFFSET $2"
	dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
	query := fmt.Sprintf("SELECT * FROM%s\nORDER BY id\n%s", dbName, params)

	var people AllPeopleResponse
	getErr := handler.repository.Get(query, &people, limit, offset)

	if getErr != nil {
		http.Error(w, getErr.Error(), http.StatusInternalServerError)
		return
	}

	response.Json(w, &people, http.StatusOK)

}

func (handler *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageParam := parseQuery(r, "page")
		sizeParam := parseQuery(r, "limit")
		if pageParam != "" && sizeParam != "" {
			handler.listLimitOffset(w, sizeParam, pageParam)
			return
		}

		if sizeParam != "" {
			handler.listLimit(w, sizeParam)
			return
		}

		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("SELECT * FROM%s\nORDER BY id", dbName)

		var people AllPeopleResponse
		dbErr := handler.repository.Get(query, &people)

		if dbErr != nil {
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, &people, http.StatusOK)

	}
}

func (handler *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
		}

		params := "WHERE id = $1"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("SELECT * FROM%s\n%s", dbName, params)

		var person Person
		dbErr := handler.repository.GetByID(query, &person, id)

		if dbErr != nil {
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, &person, http.StatusOK)

	}
}

func (handler *Handler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var person Person
		var request Request

		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&request)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		}
		errorMessages := IsValid(&request)
		if errorMessages != nil {
			http.Error(w, fmt.Sprintf("Validation failed: %s", strings.Join(errorMessages, ", ")), http.StatusBadRequest)
			return
		}

		person.Name = request.Name
		person.Surname = request.Surname
		person.Patronymic = request.Patronymic

		enrichPerson(&person)

		params := "VALUES ($1 $2 $3 $4 $5 $6)"
		returning := "RETURNING id, created_at, updated_at"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("INSERT INTO%s %s %s", dbName, params, returning)

		dbErr := handler.repository.Create(query, &person)
		if dbErr != nil {
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, &person, http.StatusCreated)

	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var person Person
		var request Request

		stringID := r.URL.Query().Get("id")
		id, parseErr := strconv.ParseUint(stringID, 10, 64)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
		}

		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&request)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		}
		errorMessages := IsValid(&request)
		if errorMessages != nil {
			http.Error(w, fmt.Sprintf("Validation failed: %s", strings.Join(errorMessages, ", ")), http.StatusBadRequest)
			return
		}

		person.Name = request.Name
		person.Surname = request.Surname
		person.Patronymic = request.Patronymic

		enrichPerson(&person)

		params := "SET name = $1, surname = $2, patronymic = $3, age = $4, gender = $5, nationality = $6)"
		selectByID := "WHERE id = $7"
		returning := "RETURNING updated_at"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("UPDATE%s\n%s\n%s\n%s", dbName, params, selectByID, returning)

		dbErr := handler.repository.Update(query, &person, id)
		if dbErr != nil {
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, &person, http.StatusOK)
	}
}

func (handler *Handler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idString := r.PathValue("id")
		id, parseErr := strconv.ParseUint(idString, 10, 64)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
		}

		params := "WHERE id = $1"
		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("DELETE FROM%s\n%s", dbName, params)

		dbErr := handler.repository.Delete(query, id)

		if dbErr != nil {
			http.Error(w, dbErr.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, nil, http.StatusNoContent)

	}
}
