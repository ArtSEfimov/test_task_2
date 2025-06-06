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
}

func (handler *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pageParam := parseQuery(r, "page")
		sizeParam := parseQuery(r, "limit")

		var page, limit uint64

		parseUint, err := strconv.ParseUint(pageParam, 10, 64)
		if err != nil && pageParam != "" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		page = parseUint

		parseUint, err = strconv.ParseUint(sizeParam, 10, 64)
		if err != nil && sizeParam != "" {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		limit = parseUint

		offset := (page - 1) * limit

		var params string
		if offset != 0 {
			params = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)
		}

		dbName := fmt.Sprintf(" %s", handler.config.Database.Name)
		query := fmt.Sprintf("SELECT * FROM%s\nORDER BY id%s", dbName, params)

		var people AllPeopleResponse

		getErr := handler.repository.GetAll(query, &people)

		if getErr != nil {
			http.Error(w, getErr.Error(), http.StatusInternalServerError)
			return
		}

		response.Json(w, &people, http.StatusOK)

	}
}

func (handler *Handler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

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

		// TODO save to Database
		createdPerson := ""

		response.Json(w, createdPerson, http.StatusCreated)

	}
}

func (handler *Handler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		stringID := r.URL.Query().Get("id")
		id, parseErr := strconv.ParseUint(stringID, 10, 64)
		if parseErr != nil {
			http.Error(w, parseErr.Error(), http.StatusBadRequest)
		}

		// TODO GET person by ID
		var person Person
		// TODO parse Database response

		var request Request
		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&request)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		}

	}
}
