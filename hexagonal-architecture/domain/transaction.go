package domain

import (
	"strings"
)

type Transaction struct {
	TransactionId   string  `db:"transaction_id"`
	AccountId       string  `db:"account_id"`
	Amount          float64 `db:"amount"`
	TransactionType string  `db:"transaction_type"`
	TransactionDate string  `db:"transaction_date"`
}

func (t Transaction) IsWithdrawal() bool {
	return strings.ToLower(t.TransactionType) == "withdrawal"
}

func (t Transaction) IsDeposit() bool {
	return strings.ToLower(t.TransactionType) == "deposit"
}
