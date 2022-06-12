package domain

import (
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/dto"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/errs"
)

type Account struct {
	AccountId   string  `db:"account_id"`
	CustomerId  string  `db:"customer_id"`
	OpeningDate string  `db:"opening_date"`
	AccountType string  `db:"account_type"`
	Amount      float64 `db:"amount"`
	Status      string  `db:"status"`
}

type AccountRepository interface {
	FindById(string) (*Account, *errs.AppError)
	Save(Account) (*Account, *errs.AppError)
	SaveTransaction(Transaction) (*Transaction, *float64, *errs.AppError)
}

func (a Account) ToNewAccountResponseDto() *dto.NewAccountResponse {
	return &dto.NewAccountResponse{
		AccountId: a.AccountId,
	}
}

func (a Account) CanWithdraw(amount float64) bool {
	return a.Amount-amount >= 0
}
