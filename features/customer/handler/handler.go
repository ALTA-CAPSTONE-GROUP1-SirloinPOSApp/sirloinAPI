package handler

import (
	"log"
	"net/http"
	"sirloinapi/features/customer"
	"sirloinapi/helper"
	"strconv"

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

func (cc *customerControl) Update() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")
		customerId := c.Param("customer_id")
		cCusId, _ := strconv.Atoi(customerId)

		updatedData := AddCustomerRequest{}
		if err := c.Bind(&updatedData); err != nil {
			return c.JSON(http.StatusBadRequest, "wrong input format")
		}

		// res, err := cc.srv.Update(token, uint(cCusId), *ToCore(input))
		_, err := cc.srv.Update(token, uint(cCusId), *ToCore(updatedData))
		if err != nil {
			return c.JSON(helper.PrintErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			// "data":    ToResponse(res),
			"message": "success edit customer data",
		})

	}
}

func (cc *customerControl) GetUserCustomers() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		res, err := cc.srv.GetUserCustomers(token)
		if err != nil {
			return c.JSON(helper.PrintErrorResponse(err.Error()))
		}
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    ToResponseArr(res),
			"message": "success get all customers",
		})
	}
}
