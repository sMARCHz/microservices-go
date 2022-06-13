package dto

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/auth-server/domain"
)

type RefreshTokenRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (r RefreshTokenRequest) IsAccessTokenValid() *jwt.ValidationError {
	_, err := jwt.Parse(r.AccessToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(domain.HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		var validationErr *jwt.ValidationError
		// Check if err is type of validationError
		if errors.As(err, &validationErr) {
			return validationErr
		}
	}
	return nil
}

func (r RefreshTokenRequest) Validate() *errs.AppError {
	if r.AccessToken == "" {
		return errs.NewAuthorizationError("missing access token")
	}
	if r.RefreshToken == "" {
		return errs.NewAuthorizationError("missing refresh token")
	}
	return nil
}
