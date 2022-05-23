package main

import (
	"encoding/json"
	"encoding/xml"
	"log"
	"net/http"
)

func main() {
	// In this section, we introduce you how to create basic api from Golang without using any third-party packages or frameworks

	// Registers the handler function for the given routes
	// In the second argument, we need to assign the function without invoking it
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/customers", getAllCustomers)

	// Listens on the TCP network address and then calls Serve with handler to handle requests on incoming connections
	// In the second argument, nil means the DefaultServeMux(ServeMux) is used
	// Http multiplexer is used to handle the incoming requests and responses
	// We can log the error message with fatal severity if the service get error
	log.Fatal(http.ListenAndServe("localhost:8080", nil))

	// We can create new mux (ServeMux)
	// mux := http.NewServeMux()
	// mux.handleFunc(.....)
	// http.ListenAndServe(...., mux)
}

type Customer struct {
	Name    string `json:"full_name" xml:"name"`
	City    string `json:"city" xml:"city"`
	Zipcode string `json:"zip_code" xml:"zipcode"`
}

func hello(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprint(w, "Hello") is ok
	w.Write([]byte("Hello"))
}

func getAllCustomers(w http.ResponseWriter, r *http.Request) {
	customers := []Customer{
		{"Anya Forger", "Westalis", "01234"},
		{"Damian Desmond", "Ostania", "56789"},
	}

	contentType := r.Header.Get("Content-Type")

	if contentType == "application/xml" {
		w.Header().Add("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(customers)
	} else {
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(customers)
	}

}
