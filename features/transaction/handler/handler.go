package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sirloinapi/features/transaction"
	"sirloinapi/helper"
	"strconv"
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
func (th *TransactionHandle) GetTransactionDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactionId := c.Param("transaction_id")
		trId, err := strconv.Atoi(transactionId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (data not found)"))
		}

		res, err := th.srv.GetTransactionDetails(uint(trId))
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
func (th *TransactionHandle) GetAdminTransactionHistory() echo.HandlerFunc {
	return func(c echo.Context) error {
		from := c.QueryParam("from")
		to := c.QueryParam("to")

		res, err := th.srv.GetAdminTransactionHistory("buy", from, to)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (data not found)"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "success get admin transaction history",
		})
	}
}
func (th *TransactionHandle) GetAdminTransactionDetails() echo.HandlerFunc {
	return func(c echo.Context) error {
		transactionId := c.Param("transaction_id")
		trId, err := strconv.Atoi(transactionId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (data not found)"))
		}

		res, err := th.srv.GetAdminTransactionDetails(uint(trId))
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (data not found)"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"data":    res,
			"message": "success get admin transaction history",
		})
	}
}
func (th *TransactionHandle) NotificationTransactionStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		// 1. Initialize empty map
		var notificationPayload map[string]interface{}

		// 2. Parse JSON request body and use it to set json to payload
		err := json.NewDecoder(c.Request().Body).Decode(&notificationPayload)
		if err != nil {
			// do something on error when decode
			return c.JSON(http.StatusBadRequest, err)
		}

		// 3. Get order-id from payload
		transactionId, exists := notificationPayload["order_id"].(string)
		if !exists {
			// do something when key `order_id` not found
			return c.JSON(http.StatusBadRequest, err)
		}

		err = th.srv.NotificationTransactionStatus(transactionId)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}

		return c.JSON(c.Response().Write([]byte("ok")))
	}
}

func (th *TransactionHandle) UpdateStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		type StatusReq struct {
			TransactionStatus string `json:"transaction_status"`
		}
		status := StatusReq{}
		err := c.Bind(&status)
		if err != nil {
			log.Println("bind order status error: ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input"))
		}

		oid := c.Param("transaction_id")
		orderId, err := strconv.Atoi(oid)
		if err != nil {
			log.Println("error read parameter: ", err.Error())
			return c.JSON(http.StatusBadRequest, helper.ErrorResponse("fail to read parameter"))
		}

		err = th.srv.UpdateStatus(uint(orderId), status.TransactionStatus)
		if err != nil {
			if strings.Contains(err.Error(), "bad request") {
				return c.JSON(http.StatusBadRequest, helper.ErrorResponse("wrong input (bad request)"))
			} else {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse("server problem"))
			}
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "success update order",
		})
	}
}
