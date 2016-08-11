package main

import (
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rivermq/rivermq/route"
	"github.com/rivermq/rivermq/util"
)

func main() {
	log.SetOutput(os.Stdout)

	router := route.NewRiverMQRouter()

	// Append PrometheusHandlerhandler
	router.
		Methods("GET").
		Path("/metrics").
		Name("PrometheusHandler").
		Handler(util.Logger(prometheus.Handler(), "PrometheusHandler"))

	log.Println("Started, listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
