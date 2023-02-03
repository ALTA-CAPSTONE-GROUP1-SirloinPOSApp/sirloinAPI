package handler

import (
	"log"
	"net/http"
	"sirloinapi/features/customer"
	"sirloinapi/helper"

	"github.com/labstack/echo/v4"
)

type customerControl struct {
	srv customer.CustomerService
}

func New(srv customer.CustomerService) customer.CustomerHandler {
	return &customerControl{
		srv: srv,
	}
}

func (cc *customerControl) Add() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		input := AddCustomerRequest{}

		if err := c.Bind(&input); err != nil {
			log.Println("error bind input")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input"))
		}

		// res, err := cc.srv.Add(token, *ToCore(input))
		_, err := cc.srv.Add(token, *ToCore(input))
		if err != nil {
			return c.JSON(helper.PrintErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			// "data":    ToResponse(res),
			"message": "success add new customer",
		})
	}
}
