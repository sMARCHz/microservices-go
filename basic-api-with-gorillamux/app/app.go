package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func Start() {
	// Create new router using gorilla/mux
	router := mux.NewRouter()

	// With gorilla/mux we can define HttpMethod of each routes
	router.HandleFunc("/customer", getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customer", createCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customer/{id:[0-9]+}", getCustomerById).Methods(http.MethodGet)

	// Implement gorialla/mux router instead of DefaultServMux
	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
