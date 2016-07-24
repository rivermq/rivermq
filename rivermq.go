package main

import (
	"log"
	"net/http"

	"github.com/rivermq/rivermq/route"
)

func main() {
	router := route.NewRiverMQRouter()
	log.Println("Started, listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
