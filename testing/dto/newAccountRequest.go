package dto

import (
	"strings"

	"github.com/sMARCHz/rest-based-microservices-go/testing/errs"
)

type NewAccountRequest struct {
	CustomerId      string  `json:"customer_id"`
	AccountType     string  `json:"account_type"`
	Amount          float64 `json:"amount"`
	TransactionType string  `json:"transaction_type"`
	TransactionDate string  `json:"transaction_date"`
}

func (r NewAccountRequest) Validate() *errs.AppError {
	if r.Amount < 5000 {
		return errs.NewValidationError("You need to deposit at least 5000.00 to open new account")
	}
	if strings.ToLower(r.AccountType) != "saving" && strings.ToLower(r.AccountType) != "checking" {
		return errs.NewValidationError("Account type should be saving or checking")
	}
	return nil
}
