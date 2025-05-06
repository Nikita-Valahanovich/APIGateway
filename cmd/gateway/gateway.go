package main

import (
	"APIGateway/pkg/handlers"
	"log"
	"net/http"
)

func main() {
	api := handlers.New()

	log.Println("API Gateway running on :8080")
	err := http.ListenAndServe(":8080", api.Router())
	if err != nil {
		log.Fatal(err)
	}
}
