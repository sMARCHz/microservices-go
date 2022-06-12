package domain

import (
	"database/sql"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/sMARCHz/rest-based-microservices-go/testing/errs"
	"github.com/sMARCHz/rest-based-microservices-go/testing/logger"
)

type AccountRepositoryDb struct {
	client *sqlx.DB
}

func (d AccountRepositoryDb) FindById(id string) (*Account, *errs.AppError) {
	findByIdSql := "SELECT account_id, customer_id, opening_date, account_type, amount, status FROM account WHERE account_id = ?"
	var account Account
	if err := d.client.Get(&account, findByIdSql, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Account not found")
		} else {
			logger.Error("Error while querying account - " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &account, nil
}

func (d AccountRepositoryDb) Save(a Account) (*Account, *errs.AppError) {
	insertSql := "INSERT INTO account (customer_id, opening_date, account_type, amount, status) VALUES(?, ?, ?, ?, ?)"
	result, err := d.client.Exec(insertSql, a.CustomerId, a.OpeningDate, a.AccountType, a.Amount, a.Status)
	if err != nil {
		logger.Error("Error while creating new account - " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	id, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting last insert id for new account - " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected error from database")
	}

	a.AccountId = strconv.FormatInt(id, 10)
	return &a, nil
}

func (d AccountRepositoryDb) SaveTransaction(t Transaction) (*Transaction, *float64, *errs.AppError) {
	tx, err := d.client.Begin()
	if err != nil {
		logger.Error("Error while starting a new transaction for bank account transaction - " + err.Error())
		return nil, nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// Insert transaction
	transactionSql := "INSERT INTO transaction (account_id, amount, transaction_type, transaction_date) VALUES(?, ?, ?, ?)"
	result, err := tx.Exec(transactionSql, t.AccountId, t.Amount, t.TransactionType, t.TransactionDate)
	if err != nil {
		logger.Error("Error while inserting transaction to the table for bank account transaction - " + err.Error())
		return nil, nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// Update balance
	if t.IsWithdrawal() {
		withDrawSql := "UPDATE account SET amount = amount - ? WHERE account_id = ?"
		_, err = tx.Exec(withDrawSql, t.Amount, t.AccountId)
	} else {
		depositSql := "UPDATE account SET amount = amount + ? WHERE account_id = ?"
		_, err = tx.Exec(depositSql, t.Amount, t.AccountId)
	}
	if err != nil {
		tx.Rollback()
		logger.Error("Error while saving transaction - " + err.Error())
		return nil, nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// Commit if everythings is good
	if err := tx.Commit(); err != nil {
		tx.Rollback()
		logger.Error("Error while committing transaction for bank account - " + err.Error())
		return nil, nil, errs.NewUnexpectedError("Unexpected database error")
	}

	transactionId, err := result.LastInsertId()
	if err != nil {
		logger.Error("Error while getting the last transaction id - " + err.Error())
		return nil, nil, errs.NewUnexpectedError("Unexpected database error")
	}
	account, appErr := d.FindById(t.AccountId)
	if appErr != nil {
		return nil, nil, appErr
	}
	t.TransactionId = strconv.FormatInt(transactionId, 10)
	return &t, &account.Amount, nil
}

func NewAccountRepositoryDb(client *sqlx.DB) AccountRepositoryDb {
	return AccountRepositoryDb{client}
}
