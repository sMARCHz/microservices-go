package app

import (
	"encoding/json"
	"net/http"

	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/dto"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/logger"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/services"
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

func (h AuthHandler) Verify(w http.ResponseWriter, r *http.Request) {
	urlParams := make(map[string]string)

	// converting from Query to map type
	for k := range r.URL.Query() {
		urlParams[k] = r.URL.Query().Get(k)
	}

	if urlParams["token"] != "" {
		isAuthorized, appError := h.service.Verify(urlParams)
		if appError != nil {
			writeResponse(w, http.StatusForbidden, unAuthorizedResponse())
		} else {
			if isAuthorized {
				writeResponse(w, http.StatusOK, authorizedResponse())
			} else {
				writeResponse(w, http.StatusForbidden, unAuthorizedResponse())
			}
		}
	} else {
		writeResponse(w, http.StatusForbidden, "missing token")
	}
}

func unAuthorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": false}
}

func authorizedResponse() map[string]bool {
	return map[string]bool{"isAuthorized": true}
}
