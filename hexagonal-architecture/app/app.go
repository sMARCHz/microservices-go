package app

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sMARCHz/microservices-go/domain"
	"github.com/sMARCHz/microservices-go/logger"
	"github.com/sMARCHz/microservices-go/services"
)

func Start() {
	sanityCheck()

	router := mux.NewRouter()

	// Wiring
	ch := CustomerHandler{services.NewCustomerService(domain.NewCustomerRepositoryDb())}

	router.HandleFunc("/customers", ch.getAllCustomers).Methods(http.MethodGet)
	router.HandleFunc("/customers", ch.createCustomer).Methods(http.MethodPost)
	router.HandleFunc("/customers/{id:[0-9]+}", ch.getCustomerById).Methods(http.MethodGet)

	address := os.Getenv("SERVER_ADDRESS")
	port := os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%v:%v", address, port), router))
}

func sanityCheck() {
	envProps := []string{
		"SERVER_ADDRESS",
		"SERVER_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_HOST",
		"DB_PORT",
		"DB_NAME",
	}
	for _, k := range envProps {
		if os.Getenv(k) == "" {
			logger.Fatal(fmt.Sprintf("Environment variable %s not defined. Terminating application...", k))
		}
	}
}
