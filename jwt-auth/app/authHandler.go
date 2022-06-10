package app

import (
	"encoding/json"
	"net/http"

	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/dto"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/logger"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/services"
)

type AuthHandler struct {
	service services.AuthService
}

func (a AuthHandler) NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "Handler not implemented...")
}

func (a AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		logger.Error("Error while decoding login request: " + err.Error())
		w.WriteHeader(http.StatusBadRequest)
	} else {
		token, err := a.service.Login(loginRequest)
		if err != nil {
			writeResponse(w, http.StatusUnauthorized, err.Error())
		} else {
			writeResponse(w, http.StatusOK, *token)
		}
	}
}

func (a AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	writeResponse(w, http.StatusOK, "Handler not implemented...")
}
