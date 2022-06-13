package domain

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go-lib/logger"
)

type AuthToken struct {
	token *jwt.Token
}

// generate access token
func (a AuthToken) GenerateAccessToken() (string, *errs.AppError) {
	signedString, err := a.token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing access token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate access token")
	}
	return signedString, nil
}

// generate refresh token from access token
func (a AuthToken) generateRefreshToken() (string, *errs.AppError) {
	claims := a.token.Claims.(AccessTokenClaims)
	refreshTokenClaims := claims.ToRefreshTokenClaims()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	signedString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing refresh token: " + err.Error())
		return "", errs.NewUnexpectedError("cannot generate refresh token")
	}
	return signedString, nil
}

// generate access token from refresh token
func GenerateAccessTokenFromRefreshToken(refreshToken string) (string, *errs.AppError) {
	token, err := jwt.ParseWithClaims(refreshToken, &RefreshTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(HMAC_SAMPLE_SECRET), nil
	})
	if err != nil {
		return "", errs.NewAuthenticationError("invalid or expired refresh token")
	}
	refreshTokenClaims := token.Claims.(*RefreshTokenClaims)
	accessTokenClaims := refreshTokenClaims.ToAccessTokenClaims()
	authToken := NewAuthToken(accessTokenClaims)
	return authToken.GenerateAccessToken()
}

func NewAuthToken(claims AccessTokenClaims) AuthToken {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return AuthToken{token: token}
}
