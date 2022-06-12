package domain

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go-lib/logger"
)

type AuthRepository interface {
	ValidateUser(string, string) (*UserDetail, *errs.AppError)
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (a AuthRepositoryDb) ValidateUser(username, password string) (*UserDetail, *errs.AppError) {
	var user UserDetail
	validateUserSql := `SELECT username, u.customer_id, role, GROUP_CONCAT(a.account_id) AS account_numbers FROM user u
					LEFT JOIN account a
					ON a.customer_id = u.customer_id
					WHERE username = ? and password = ?
					GROUP BY username, customer_id, role`
	if err := a.client.Get(&user, validateUserSql, username, password); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}
	return &user, nil
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
