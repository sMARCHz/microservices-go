package domain

type CustomerRepositoryStub struct {
	customers []Customer
}

func (c CustomerRepositoryStub) FindAll() ([]Customer, error) {
	return c.customers, nil
}

func NewCustomerRepositoryStub() CustomerRepositoryStub {
	customers := []Customer{
		{Id: "1", Name: "Loid Forger", City: "Westalis", Zipcode: "01234"},
		{Id: "2", Name: "Yor Forger", City: "Ostania", Zipcode: "01234"},
	}
	return CustomerRepositoryStub{customers: customers}
}
