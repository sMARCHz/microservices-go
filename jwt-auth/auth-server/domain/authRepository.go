package domain

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go-lib/logger"
)

type AuthRepository interface {
	ValidateUser(string, string) (*Login, *errs.AppError)
	GenerateAndSaveRefreshToken(AuthToken) (string, *errs.AppError)
	IsRefreshTokenExisted(string) *errs.AppError
}

type AuthRepositoryDb struct {
	client *sqlx.DB
}

func (a AuthRepositoryDb) ValidateUser(username, password string) (*Login, *errs.AppError) {
	var login Login
	validateUserSql := `SELECT username, u.customer_id, role, GROUP_CONCAT(a.account_id) AS account_numbers FROM user u
					LEFT JOIN account a
					ON a.customer_id = u.customer_id
					WHERE username = ? and password = ?
					GROUP BY username, customer_id, role`
	if err := a.client.Get(&login, validateUserSql, username, password); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewAuthenticationError("invalid credentials")
		} else {
			logger.Error("Error while verifying login request from database: " + err.Error())
			return nil, errs.NewUnexpectedError("unexpected database error")
		}
	}
	return &login, nil
}

func (a AuthRepositoryDb) GenerateAndSaveRefreshToken(authToken AuthToken) (string, *errs.AppError) {
	// Generate refresh token
	var refreshToken string
	var appErr *errs.AppError
	if refreshToken, appErr = authToken.generateRefreshToken(); appErr != nil {
		return "", appErr
	}
	// Save refresh token
	insertRefreshTokenSql := "INSERT INTO refresh_token_store (refresh_token) VALUES (?)"
	if _, err := a.client.Exec(insertRefreshTokenSql, refreshToken); err != nil {
		logger.Error("unexpected database error: " + err.Error())
		return "", errs.NewUnexpectedError("unexpected database error")
	}
	return refreshToken, nil
}

func (a AuthRepositoryDb) IsRefreshTokenExisted(refreshToken string) *errs.AppError {
	selectSql := "SELECT refresh_token FROM refresh_token_store WHERE refresh_token = ?"
	var token string
	err := a.client.Get(&token, selectSql, refreshToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return errs.NewAuthenticationError("refresh token isn't registered in the store")
		} else {
			logger.Error("unexpected database error: " + err.Error())
			return errs.NewUnexpectedError("unexpected database error")
		}
	}
	return nil
}

func NewAuthRepository(client *sqlx.DB) AuthRepositoryDb {
	return AuthRepositoryDb{client}
}
