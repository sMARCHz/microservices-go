package services

import (
	"time"

	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/domain"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/dto"
)

type AccountService interface {
	CreateAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	MakeTransaction(dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}

type DefaultAccountService struct {
	repo domain.AccountRepository
}

func (a DefaultAccountService) CreateAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	account := domain.Account{
		CustomerId:  req.CustomerId,
		OpeningDate: time.Now().Format("2006-01-02 15:04:05"),
		AccountType: req.AccountType,
		Amount:      req.Amount,
		Status:      "1",
	}
	if newAccount, err := a.repo.Save(account); err != nil {
		return nil, err
	} else {
		return newAccount.ToNewAccountResponseDto(), nil
	}
}

func (a DefaultAccountService) MakeTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if req.IsTransactionTypeWithdrawal() {
		account, err := a.repo.FindById(req.AccountId)
		if err != nil {
			return nil, err
		}
		if !account.CanWithdraw(req.Amount) {
			return nil, errs.NewValidationError("Insufficient balance in the account")
		}
	}

	t := domain.Transaction{
		AccountId:       req.AccountId,
		Amount:          req.Amount,
		TransactionType: req.TransactionType,
		TransactionDate: time.Now().Format("2006-01-02 15:04:05"),
	}

	if transaction, balance, err := a.repo.SaveTransaction(t); err != nil {
		return nil, err
	} else {
		return transaction.ToResponseDto(*balance), nil
	}
}

func NewAccountService(repository domain.AccountRepository) DefaultAccountService {
	return DefaultAccountService{repository}
}
