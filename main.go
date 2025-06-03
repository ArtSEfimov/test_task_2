package main

import (
	"fmt"
	"go_test_task_2/config"
	"log"
	"net/http"
)

const Port = ":8080" // TODO relocate it to .env

func main() {

	// load config
	conf := config.NewConfig()

	// init mux
	peopleMux := http.NewServeMux()

	// create new server
	peopleServer := http.Server{
		Addr:    Port,
		Handler: peopleMux,
	}

	bytes := make([]byte, 1024)
	r, _ := http.Get("https://api.agify.io?name=ivan")
	n, _ := r.Body.Read(bytes)
	fmt.Println(string(bytes[:n]))

	// TODO create new Repo
	// TODO create DB

	conf.InfoLogger.Printf("Starting server on port %s...", Port)
	serverStartErr := peopleServer.ListenAndServe()
	if serverStartErr != nil {
		log.Fatal(serverStartErr)
	}

}
