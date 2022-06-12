package domain

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/sMARCHz/rest-based-microservices-go-lib/errs"
	"github.com/sMARCHz/rest-based-microservices-go-lib/logger"
)

type CustomerRepositoryDb struct {
	client *sqlx.DB
}

func (c CustomerRepositoryDb) FindAll() ([]Customer, *errs.AppError) {
	findAllSql := "SELECT customer_id, name, city, zipcode, date_of_birth, status FROM customer"
	customers := make([]Customer, 0)

	// If the fields don't match with database column, we need to define db tag in the domain
	err := c.client.Select(&customers, findAllSql)
	if err != nil {
		logger.Error("Error while querying customer - " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}

	// We can use structScan instead of scan (better if use Select() in line20)
	// err = sqlx.StructScan(rows, &customers)
	// if err != nil {
	// 	logger.Error("Error while scanning customers - " + err.Error())
	// 	return nil, errs.NewUnexpectedError("Unexpected database error")
	// }
	return customers, nil
}

func (c CustomerRepositoryDb) FindById(id string) (*Customer, *errs.AppError) {
	findByIdSql := "SELECT customer_id, name, city, zipcode, date_of_birth, status FROM customer WHERE customer_id = ?"
	var customer Customer

	err := c.client.Get(&customer, findByIdSql, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errs.NewNotFoundError("Customer not found")
		} else {
			logger.Error("Error while querying customer - " + err.Error())
			return nil, errs.NewUnexpectedError("Unexpected database error")
		}
	}
	return &customer, nil
}

func (c CustomerRepositoryDb) FindByStatus(status string) ([]Customer, *errs.AppError) {
	findByStatusSql := "SELECT customer_id, name, city, zipcode, date_of_birth, status FROM customer WHERE status = ?"
	customers := make([]Customer, 0)

	err := c.client.Select(&customers, findByStatusSql, status)
	if err != nil {
		logger.Error("Error while querying customer - " + err.Error())
		return nil, errs.NewUnexpectedError("Unexpected database error")
	}
	return customers, nil
}

func NewCustomerRepositoryDb(client *sqlx.DB) CustomerRepositoryDb {
	return CustomerRepositoryDb{client}
}
