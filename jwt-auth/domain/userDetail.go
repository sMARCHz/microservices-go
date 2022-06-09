package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sMARCHz/microservices-go/jwt-auth/logger"
)

const TOKEN_DURATION = time.Hour
const HMAC_SAMPLE_SECRET = "hmacSampleSecret"

type UserDetail struct {
	Username   string `db:"username"`
	CustomerId string `db:"customer_id"`
	Accounts   string `db:"account_numbers"`
	Role       string `db:"role"`
}

func (u UserDetail) GenerateToken() (*string, error) {
	var claims jwt.MapClaims
	if u.Role == "admin" {
		claims = u.claimsForAdmin()
	} else {
		claims = u.claimsForUser()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedTokenAsString, err := token.SignedString([]byte(HMAC_SAMPLE_SECRET))
	if err != nil {
		logger.Error("Failed while signing token: " + err.Error())
		return nil, errors.New("cannot generate token")
	}
	return &signedTokenAsString, nil
}

func (u UserDetail) claimsForUser() jwt.MapClaims {
	accounts := strings.Split(u.Accounts, ",")
	return jwt.MapClaims{
		"customer_id": u.CustomerId,
		"role":        u.Role,
		"username":    u.Username,
		"accounts":    accounts,
		"exp":         time.Now().Add(TOKEN_DURATION).Unix(),
	}
}

func (u UserDetail) claimsForAdmin() jwt.MapClaims {
	return jwt.MapClaims{
		"role":     u.Role,
		"username": u.Username,
		"exp":      time.Now().Add(TOKEN_DURATION).Unix(),
	}
}
