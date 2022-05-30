package app

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sMARCHz/microservices-go/domain"
	"github.com/sMARCHz/microservices-go/services"
)

func Start() {
	router := mux.NewRouter()

	// Wiring
	ch := CustomerHandler{services.NewCustomerService(domain.NewCustomerRepositoryDb())}

	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers", ch.createCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customers/{id:[0-9]+}", ch.getCustomerById).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe("localhost:8080", router))
}
