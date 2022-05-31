package services

import (
	"time"

	"github.com/sMARCHz/microservices-go/domain"
	"github.com/sMARCHz/microservices-go/dto"
	"github.com/sMARCHz/microservices-go/errs"
)

type AccountService interface {
	CreateAccount(dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError)
	MakeTransaction(dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError)
}

type AccountServiceImpl struct {
	repo domain.AccountRepository
}

func (a AccountServiceImpl) CreateAccount(req dto.NewAccountRequest) (*dto.NewAccountResponse, *errs.AppError) {
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

func (a AccountServiceImpl) MakeTransaction(req dto.TransactionRequest) (*dto.TransactionResponse, *errs.AppError) {
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

	transaction, balance, err := a.repo.SaveTransaction(t)
	if err != nil {
		return nil, err
	}
	return &dto.TransactionResponse{TransactionId: transaction.TransactionId, Balance: *balance}, nil
}

func NewAccountService(repository domain.AccountRepository) AccountServiceImpl {
	return AccountServiceImpl{repository}
}
