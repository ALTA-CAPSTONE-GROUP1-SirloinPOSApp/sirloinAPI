package handler

import (
	"log"
	"net/http"
	"sirloinapi/features/transaction"
	"sirloinapi/helper"
	"strings"

	"github.com/labstack/echo/v4"
)

type TransactionHandle struct {
	srv transaction.TransactionService
}

func New(ts transaction.TransactionService) transaction.TransactionHandler {
	return &TransactionHandle{
		srv: ts,
	}
}

func (th *TransactionHandle) AddSell() echo.HandlerFunc {
	return func(c echo.Context) error {
		uCart := transaction.Cart{}
		if err := c.Bind(&uCart); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (bad request)"))
		}

		token := c.Get("user")
		if token == nil {
			log.Println("error get token JWT")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (bad request)"))
		}

		res, err := th.srv.AddSell(token, uCart)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") || strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (bad request)"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    res,
			"message": "transaction created",
		})
	}
}
func (th *TransactionHandle) AddBuy() echo.HandlerFunc {
	return func(c echo.Context) error {
		uCart := transaction.Cart{}
		if err := c.Bind(&uCart); err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (bad request)"))
		}

		token := c.Get("user")
		if token == nil {
			log.Println("error get token JWT")
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (bad request)"))
		}

		res, err := th.srv.AddBuy(token, uCart)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") || strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (bad request)"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"data":    res,
			"message": "transaction created",
		})
	}
}
func (th *TransactionHandle) GetTransactionHistory() echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Get("user")

		from := c.QueryParam("from")
		to := c.QueryParam("to")
		status := c.QueryParam("status")

		res, err := th.srv.GetTransactionHistory(token, status, from, to)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (data not found)"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "success get transaction history",
		})
	}
}
