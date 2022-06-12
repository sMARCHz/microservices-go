package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go-lib/logger"
)

type AuthToken struct {
	token *jwt.Token
}

func (a AuthToken) NewAccessToken() (string, *errs.AppError) {
	signedString, err := a.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing access token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate token")
	}
	return signedString, nil
}

func NewAuthToken(claims Claims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}
