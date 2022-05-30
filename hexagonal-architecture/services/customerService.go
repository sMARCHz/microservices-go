package services

import (
	"github.com/sMARCHz/microservices-go/domain"
	"github.com/sMARCHz/microservices-go/errs"
)

type CustomerService interface {
	GetAllCustomers(string) ([]domain.Customer, *errs.AppError)
	GetCustomerById(string) (*domain.Customer, *errs.AppError)
}

type CustomerServiceImpl struct {
	repo domain.CustomerRepository
}

func (c CustomerServiceImpl) GetAllCustomers(status string) ([]domain.Customer, *errs.AppError) {
	if status == "active" {
		return c.repo.FindByStatus("1")
	} else if status == "inactive" {
		return c.repo.FindByStatus("0")
	} else {
		return c.repo.FindAll()
	}
}

func (c CustomerServiceImpl) GetCustomerById(id string) (*domain.Customer, *errs.AppError) {
	return c.repo.FindById(id)
}

func NewCustomerService(repository domain.CustomerRepository) CustomerServiceImpl {
	return CustomerServiceImpl{repository}
}
