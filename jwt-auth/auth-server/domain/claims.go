package domain

import (
	"os"

	"github.com/dgrijalva/jwt-go"
)

var HMAC_SAMPLE_SECRET = os.Getenv("JWT_SECRET")

type Claims struct {
	CustomerId string   `json:"customer_id"`
	Accounts   []string `json:"accounts"`
	Username   string   `json:"username"`
	Role       string   `json:"role"`
	jwt.StandardClaims
}

func (c Claims) IsUserRole() bool {
	return c.Role == "user"
}

func (c Claims) IsValidCustomerId(customerId string) bool {
	return c.CustomerId == customerId
}

func (c Claims) IsValidAccountId(accountId string) bool {
	if accountId != "" {
		accountFound := false
		for _, a := range c.Accounts {
			if a == accountId {
				accountFound = true
				break
			}
		}
		return accountFound
	}
	return true
}

func (c Claims) IsRequestVerifiedWithTokenClaims(urlParams map[string]string) bool {
	if !c.IsValidCustomerId(urlParams["customer_id"]) {
		return false
	}

	if !c.IsValidAccountId(urlParams["account_id"]) {
		return false
	}
	return true
}
