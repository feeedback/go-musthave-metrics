package main

import (
	"net/http"

	"github.com/feeedback/go-musthave-metrics/internal/handlers"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	return http.ListenAndServe(`:8080`, http.HandlerFunc(handlers.Webhook))
}
