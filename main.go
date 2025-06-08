package main

import (
	"fmt"
	"go_test_task_2/config"
	"go_test_task_2/internal/people"
	"go_test_task_2/pkg/db"
	"log"
	"net/http"
)

func main() {

	// load config
	conf := config.NewConfig()
	port := conf.Port

	// init mux
	peopleMux := http.NewServeMux()

	// create new server
	peopleServer := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: peopleMux,
	}

	peopleDB := db.NewDB(conf)
	peopleRepository := people.NewRepository(peopleDB)

	people.NewHandler(peopleMux, people.NewHandlerDeps(conf, peopleRepository))

	conf.InfoLogger.Printf("Starting server on port %s...", port)
	serverStartErr := peopleServer.ListenAndServe()
	if serverStartErr != nil {
		log.Fatal(serverStartErr)
	}

}
