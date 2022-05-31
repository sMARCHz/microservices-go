package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sMARCHz/microservices-go/services"
)

type CustomerHandler struct {
	service services.CustomerService
}

func (c CustomerHandler) getAllCustomers(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	customers, err := c.service.GetAllCustomers(status)
	if err != nil {
		writeResponse(w, err.Code, err)
	} else {
		writeResponse(w, http.StatusOK, customers)
	}
}

func (c CustomerHandler) createCustomer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Post request received")
}

func (c CustomerHandler) getCustomerById(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	customer, err := c.service.GetCustomerById(id)
	if err != nil {
		writeResponse(w, err.Code, err)
	} else {
		writeResponse(w, http.StatusOK, customer)
	}
}

func writeResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
