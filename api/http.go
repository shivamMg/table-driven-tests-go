package api

import (
	"log"
	"net/http"
)

func respond(w http.ResponseWriter, statusCode int, response string) {
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(response)); err != nil {
		log.Println("failed to write response:", err)
	}
}
