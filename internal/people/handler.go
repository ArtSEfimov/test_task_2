package people

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
}

func NewHandler(router *http.ServeMux) {
	handler := &Handler{}
	router.HandleFunc("GET /people", handler.GetAll())
	router.HandleFunc("GET /people/{id}", handler.Get())
	router.HandleFunc("POST /people", handler.Create())
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

		var enrichErr error
		enrichErr = enrichPerson(&person)
		if enrichErr != nil {
			log.Println(enrichErr)
		}

	}
}
