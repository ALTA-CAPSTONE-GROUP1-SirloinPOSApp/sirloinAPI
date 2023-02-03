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
}
type CustomerService interface {
	Add(userToken interface{}, newCustomer Core) (Core, error)
}
type CustomerData interface {
	Add(userId uint, newCustomer Core) (Core, error)
}
