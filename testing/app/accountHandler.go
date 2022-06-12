package app

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sMARCHz/rest-based-microservices-go/testing/dto"
	"github.com/sMARCHz/rest-based-microservices-go/testing/services"
)

type AccountHandler struct {
	service services.AccountService
}

func (a AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) {
	var req dto.NewAccountRequest
	var customerId = mux.Vars(r)["customer_id"]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		req.CustomerId = customerId
		if account, appErr := a.service.CreateAccount(req); appErr != nil {
			writeResponse(w, appErr.Code, appErr)
		} else {
			writeResponse(w, http.StatusCreated, account)
		}
	}
}

func (a AccountHandler) makeTransaction(w http.ResponseWriter, r *http.Request) {
	var req dto.TransactionRequest
	var customerId = mux.Vars(r)["customer_id"]
	var accountId = mux.Vars(r)["account_id"]
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeResponse(w, http.StatusBadRequest, err.Error())
	} else {
		req.CustomerId = customerId
		req.AccountId = accountId
		if response, appErr := a.service.MakeTransaction(req); appErr != nil {
			writeResponse(w, appErr.Code, appErr)
		} else {
			writeResponse(w, http.StatusOK, response)
		}
	}
}
