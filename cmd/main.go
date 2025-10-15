package main

import (
	"log"
	"net/http"

	"github.com/ValeriyL01/balance-service/internal/handlers"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/user", handlers.GetUserBalance)

	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)
}
