package people

import (
	"bufio"
	"encoding/json"
	"go_test_task_2/config"
	"go_test_task_2/pkg/response"
	"net/http"
	"strconv"
)

type HandlerDeps struct {
	Config     *config.Config
	Repository *Repository
}
type Handler struct {
}

func NewHandler(router *http.ServeMux) {
	handler := &Handler{}
	router.HandleFunc("GET /people", handler.GetAll())
	router.HandleFunc("GET /people/{id}", handler.Get())
	router.HandleFunc("POST /people", handler.Create())
	router.HandleFunc("PUT  /people/{id}", handler.Update())
}

func (handler *Handler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}

func (handler *Handler) Get() http.HandlerFunc {
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

		person.Name = request.Name
		person.Surname = request.Surname
		person.Patronymic = request.Patronymic

		enrichPerson(&person)

		// TODO save to DB
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
		// TODO parse DB response

		var request Request
		bodyReader := bufio.NewReader(r.Body)
		decodeErr := json.NewDecoder(bodyReader).Decode(&request)
		if decodeErr != nil {
			http.Error(w, decodeErr.Error(), http.StatusBadRequest)
		}

	}
}
