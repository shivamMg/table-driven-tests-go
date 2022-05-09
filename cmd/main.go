package main

import (
	"log"
	"net/http"

	"github.com/shivamMg/table-driven-tests-go/api"
)

func main() {
	ctrl := api.NewController(&api.AuthClient{}, &api.DBClient{})
	http.HandleFunc("/todos", ctrl.CreateTODO)

	hostPort := "localhost:8080"
	log.Println("listening on", hostPort)
	if err := http.ListenAndServe(hostPort, nil); err != nil {
		log.Fatal(err)
	}
}
