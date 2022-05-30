package domain

import "github.com/sMARCHz/microservices-go/errs"

type Customer struct {
	Id          string `json:"customer_id"`
	Name        string `json:"full_name"`
	City        string `json:"city"`
	Zipcode     string `json:"zip_code"`
	DateofBirth string `json:"date_of_birth"`
	Status      string `json:"status"`
}

type CustomerRepository interface {
	FindAll() ([]Customer, *errs.AppError)
	FindById(string) (*Customer, *errs.AppError)
	FindByStatus(string) ([]Customer, *errs.AppError)
}
