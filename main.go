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
		Addr:    port,
		Handler: peopleMux,
	}

	bytes := make([]byte, 1024)
	r, _ := http.Get("https://api.agify.io?name=ivan")
	n, _ := r.Body.Read(bytes)
	fmt.Println(string(bytes[:n]))

	peopleDB := db.NewDB(conf)
	peopleRepository := people.NewRepository(peopleDB)
	fmt.Println(peopleRepository) // TODO to handler

	conf.InfoLogger.Printf("Starting server on port %s...", port)
	serverStartErr := peopleServer.ListenAndServe()
	if serverStartErr != nil {
		log.Fatal(serverStartErr)
	}

}
