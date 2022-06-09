package domain

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/sMARCHz/microservices-go/jwt-auth/logger"
)

type AuthRepository interface {
	ValidateUser(string, string) (*UserDetail, error)
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (a AuthRepositoryDb) ValidateUser(username, password string) (*UserDetail, error) {
	var user UserDetail
	validateUserSql := `SELECT username, u.customer_id, role, GROUP_CONCAT(a.account_id) AS account_numbers FROM user u
					LEFT JOIN account a
					ON a.customer_id = u.customer_id
					WHERE username = ? and password = ?
					GROUP BY username, customer_id, role`
	if err := a.client.Get(&user, validateUserSql, username, password); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errors.New("unexpected database error")
		}
	}
	return &user, nil
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
