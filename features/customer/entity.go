package customer

import "github.com/labstack/echo/v4"

type Core struct {
	ID          uint
	Email       string
	Name        string
	PhoneNumber string
	Address     string
}

type CustomerHandler interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	GetUserCustomers() echo.HandlerFunc
}
type CustomerService interface {
	Add(userToken interface{}, newCustomer Core) (Core, error)
	Update(userToken interface{}, customerId uint, updateData Core) (Core, error)
	GetUserCustomers(token interface{}) ([]Core, error)
}
type CustomerData interface {
	Add(userId uint, newCustomer Core) (Core, error)
	Update(userId, customerId uint, updateData Core) (Core, error)
	GetUserCustomers(userId uint) ([]Core, error)
}
