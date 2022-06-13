package dto

import "github.com/sMARCHz/rest-based-microservices-go-lib/errs"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (l LoginRequest) Validate() *errs.AppError {
	if l.Username == "" {
		return errs.NewAuthenticationError("missing username")
	}
	if l.Password == "" {
		return errs.NewAuthenticationError("missing password")
	}
	return nil
}
