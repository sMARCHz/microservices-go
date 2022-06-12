package dto

import (
	"strings"

	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/errs"
)

type TransactionRequest struct {
	CustomerId      string  `json:"customer_id"`
	AccountId       string  `json:"account_id"`
	TransactionType string  `json:"transaction_type"`
	Amount          float64 `json:"amount"`
}

func (t TransactionRequest) Validate() *errs.AppError {
	if !t.IsTransactionTypeDeposit() && !t.IsTransactionTypeWithdrawal() {
		return errs.NewValidationError("Transaction type can only be deposit or withdrawal")
	}
	if t.Amount <= 0 {
		return errs.NewValidationError("Amount cannot be less than or equal to zero")
	}
	return nil
}

func (t TransactionRequest) IsTransactionTypeWithdrawal() bool {
	return strings.ToLower(t.TransactionType) == "withdrawal"
}

func (t TransactionRequest) IsTransactionTypeDeposit() bool {
	return strings.ToLower(t.TransactionType) == "deposit"
}
