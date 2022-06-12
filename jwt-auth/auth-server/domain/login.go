package domain

import (
	"database/sql"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const TOKEN_DURATION = time.Hour

type Login struct {
	Username   string         `db:"username"`
	CustomerId sql.NullString `db:"customer_id"`
	Accounts   sql.NullString `db:"account_numbers"`
	Role       string         `db:"role"`
}

func (l Login) ClaimsForAccessToken() Claims {
	if l.Role == "admin" {
		return l.claimsForAdmin()
	} else {
		return l.claimsForUser()
	}
}

func (l Login) claimsForUser() Claims {
	accounts := strings.Split(l.Accounts.String, ",")
	return Claims{
		CustomerId: l.CustomerId.String,
		Accounts:   accounts,
		Username:   l.Username,
		Role:       l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TOKEN_DURATION).Unix(),
		},
	}
}

func (l Login) claimsForAdmin() Claims {
	return Claims{
		Username: l.Username,
		Role:     l.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(TOKEN_DURATION).Unix(),
		},
	}
}
