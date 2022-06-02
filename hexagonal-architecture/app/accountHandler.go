package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sMARCHz/microservices-go/hexagonal-architecture/dto"
	"github.com/sMARCHz/microservices-go/hexagonal-architecture/services"
)

type AccountHandler struct {
	service services.AccountService
}

func (a AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	var request dto.NewAccountRequest
	var customerId = mux.Vars(r)["customer_id"]
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		request.CustomerId = customerId
		if account, appError := a.service.CreateAccount(request); appError != nil {
			writeResponse(w, appError.Code, appError)
		} else {
			writeResponse(w, http.StatusCreated, account)
		}
	}
}

func (a AccountHandler) MakeTransaction(w http.ResponseWriter, r *http.Request) {
	var request dto.TransactionRequest
	var customerId = mux.Vars(r)["customer_id"]
	var accountId = mux.Vars(r)["account_id"]
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		request.CustomerId = customerId
		request.AccountId = accountId
		response, err := a.service.MakeTransaction(request)
		if err != nil {
			writeResponse(w, err.Code, err)
		} else {
			writeResponse(w, http.StatusOK, response)
		}
	}
}
