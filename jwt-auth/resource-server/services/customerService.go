package services

import (
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/domain"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/dto"
	"github.com/sMARCHz/rest-based-microservices-go/jwt-auth/resource-server/errs"
)

type CustomerService interface {
	GetAllCustomers(string) ([]dto.CustomerResponse, *errs.AppError)
	GetCustomerById(string) (*dto.CustomerResponse, *errs.AppError)
}

type DefaultCustomerService struct {
	repo domain.CustomerRepository
}

func (c DefaultCustomerService) GetAllCustomers(status string) ([]dto.CustomerResponse, *errs.AppError) {
	var customers []domain.Customer
	var err *errs.AppError
	if status == "active" {
		customers, err = c.repo.FindByStatus("1")
	} else if status == "inactive" {
		customers, err = c.repo.FindByStatus("0")
	} else {
		customers, err = c.repo.FindAll()
	}

	if err != nil {
		return nil, err
	}
	response := make([]dto.CustomerResponse, 0)
	for _, v := range customers {
		response = append(response, v.ToDto())
	}
	return response, nil
}

func (c DefaultCustomerService) GetCustomerById(id string) (*dto.CustomerResponse, *errs.AppError) {
	customer, err := c.repo.FindById(id)
	if err != nil {
		return nil, err
	}
	response := customer.ToDto()
	return &response, nil
}

func NewCustomerService(repository domain.CustomerRepository) DefaultCustomerService {
	return DefaultCustomerService{repository}
}
