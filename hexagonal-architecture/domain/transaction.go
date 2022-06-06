package domain

import (
	"strings"

	"github.com/sMARCHz/microservices-go/hexagonal-architecture/dto"
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

func (t Transaction) ToResponseDto(balance float64) *dto.TransactionResponse {
	return &dto.TransactionResponse{
		TransactionId: t.TransactionId,
		Balance:       balance,
	}
}
