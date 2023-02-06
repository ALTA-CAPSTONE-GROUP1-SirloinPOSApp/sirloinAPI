package services

import (
	"errors"
	"sirloinapi/config"
	"sirloinapi/features/transaction"
	"sirloinapi/helper"
	"sirloinapi/mocks"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestAddSell(t *testing.T) {
	data := mocks.NewTransactionData(t)
	userId := 1
	newCart := transaction.Cart{
		Items: []transaction.Item{
			{
				ProductId: 1,
				Quantity:  15,
				Price:     5000,
			},
			{
				ProductId: 2,
				Quantity:  10,
				Price:     20000,
			},
			{
				ProductId: 7,
				Quantity:  20,
				Price:     12000,
			},
		},
	}
	expectedData := transaction.Core{
		ID:                1,
		CustomerId:        1,
		TotalPrice:        550000,
		TransactionStatus: "pending",
		PaymentUrl:        "qris.png",
	}
	t.Run("transaction created", func(t *testing.T) {
		data.On("AddSell", uint(userId), newCart).Return(expectedData, nil).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.AddSell(pToken, newCart)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		assert.Equal(t, expectedData.CustomerId, res.CustomerId)
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)

		res, err := srv.AddSell(token, newCart)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("AddSell", uint(userId), newCart).Return(transaction.Core{}, errors.New("server problem")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.AddSell(pToken, newCart)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		data.On("AddSell", uint(userId), newCart).Return(transaction.Core{}, errors.New("bad request")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.AddSell(pToken, newCart)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})
}

func ToTime(t string) time.Time {
	res, _ := time.Parse("2006-01-02 15:04:05", t)
	return res
}

func TestAddBuy(t *testing.T) {
	data := mocks.NewTransactionData(t)
	userId := 1
	newCart := transaction.Cart{
		Items: []transaction.Item{
			{
				ProductId: 1,
				Quantity:  15,
				Price:     5000,
			},
			{
				ProductId: 2,
				Quantity:  10,
				Price:     20000,
			},
			{
				ProductId: 7,
				Quantity:  20,
				Price:     12000,
			},
		},
	}
	expectedData := transaction.Core{
		ID:                1,
		CustomerId:        1,
		TotalPrice:        550000,
		TransactionStatus: "pending",
		PaymentUrl:        "qris.png",
	}
	t.Run("transaction created", func(t *testing.T) {
		data.On("AddBuy", uint(userId), newCart).Return(expectedData, nil).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.AddBuy(pToken, newCart)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		assert.Equal(t, expectedData.CustomerId, res.CustomerId)
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)

		res, err := srv.AddBuy(token, newCart)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("AddBuy", uint(userId), newCart).Return(transaction.Core{}, errors.New("server problem")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.AddBuy(pToken, newCart)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})

	t.Run("bad request", func(t *testing.T) {
		data.On("AddBuy", uint(userId), newCart).Return(transaction.Core{}, errors.New("bad request")).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.AddBuy(pToken, newCart)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "bad request")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})
}

func TestGetTransactionHistory(t *testing.T) {
	data := mocks.NewTransactionData(t)
	userId := 2
	sendEmail := "true"
	expectedData := []transaction.Core{
		{
			ID:                1,
			CustomerId:        1,
			CustomerName:      "customer1",
			TotalPrice:        550000,
			Discount:          0.1,
			TotalBill:         495000,
			CreatedAt:         ToTime("2023-01-26T02:11:48"),
			TransactionStatus: "success",
			InvoiceNumber:     "INV/01",
			InvoiceUrl:        "invoice.pdf",
			PaymentUrl:        "url",
		},
		{
			ID:                2,
			CustomerId:        2,
			CustomerName:      "customer2",
			TotalPrice:        1000000,
			Discount:          0.1,
			TotalBill:         900000,
			CreatedAt:         ToTime("2023-01-27T02:11:48"),
			TransactionStatus: "success",
			InvoiceNumber:     "INV/02",
			InvoiceUrl:        "invoice.pdf",
			PaymentUrl:        "url",
		},
	}
	from := "2022-01-01"
	to := "2022-12-31"
	status := "sell"
	t.Run("success get transaction history", func(t *testing.T) {
		data.On("GetTransactionHistory", uint(userId), status, from, to, sendEmail).Return(expectedData, nil).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetTransactionHistory(pToken, status, from, to, sendEmail)
		assert.Nil(t, err)
		assert.Equal(t, len(res), len(expectedData))
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("GetTransactionHistory", uint(userId), status, from, to, sendEmail).Return([]transaction.Core{}, errors.New("server problem")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetTransactionHistory(pToken, status, from, to, sendEmail)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, 0, len(res))
		data.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("GetTransactionHistory", uint(userId), status, from, to, sendEmail).Return([]transaction.Core{}, errors.New("data not found")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetTransactionHistory(pToken, status, from, to, sendEmail)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, 0, len(res))
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(1)

		res, err := srv.GetTransactionHistory(token, status, from, to, sendEmail)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user not found")
		assert.Equal(t, 0, len(res))
	})

}

func TestGetTransactionDetails(t *testing.T) {
	data := mocks.NewTransactionData(t)
	transactionId := 1
	expectedData := transaction.TransactionRes{
		ID:                1,
		CustomerId:        1,
		CustomerName:      "customer1",
		TotalPrice:        550000,
		Discount:          0.1,
		TotalBill:         495000,
		CreatedAt:         ToTime("2023-01-26T02:11:48"),
		TransactionStatus: "success",
		InvoiceUrl:        "https://mediasosial.s3.ap-southeast-1.amazonaws.com/invoice/InvoiceSimple-PDF-Template.pdf",
		TransactionProductRes: []transaction.TransactionProductRes{
			{
				ProductId: 1,
				Quantity:  15,
				Price:     5000,
			},
			{
				ProductId: 2,
				Quantity:  10,
				Price:     20000,
			},
			{
				ProductId: 7,
				Quantity:  20,
				Price:     12000,
			},
		},
	}

	t.Run("success get transaction history", func(t *testing.T) {
		data.On("GetTransactionDetails", uint(transactionId)).Return(expectedData, nil).Once()
		srv := New(data)

		res, err := srv.GetTransactionDetails(uint(transactionId))
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("GetTransactionDetails", uint(transactionId)).Return(transaction.TransactionRes{}, errors.New("server problem")).Once()
		srv := New(data)

		res, err := srv.GetTransactionDetails(uint(transactionId))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})
	t.Run("data not found", func(t *testing.T) {
		data.On("GetTransactionDetails", uint(transactionId)).Return(transaction.TransactionRes{}, errors.New("data not found")).Once()
		srv := New(data)

		res, err := srv.GetTransactionDetails(uint(transactionId))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})
}

func TestGetAdminTransactionHistory(t *testing.T) {
	data := mocks.NewTransactionData(t)
	userId := 2

	expectedData := []transaction.AdmTransactionRes{
		{
			ID:                1,
			TenantId:          1,
			TenantName:        "customer1",
			TotalBill:         495000,
			CreatedAt:         ToTime("2023-01-26T02:11:48"),
			TransactionStatus: "success",
			InvoiceNumber:     "INV/01",
			InvoiceUrl:        "invoice.pdf",
			PaymentUrl:        "url",
		},
		{
			ID:                2,
			TenantId:          2,
			TenantName:        "customer2",
			TotalBill:         900000,
			CreatedAt:         ToTime("2023-01-27T02:11:48"),
			TransactionStatus: "success",
			InvoiceNumber:     "INV/02",
			InvoiceUrl:        "invoice.pdf",
			PaymentUrl:        "url",
		},
	}
	from := "2022-01-01"
	to := "2022-12-31"
	status := "sell"
	t.Run("success get admin transaction history", func(t *testing.T) {
		data.On("GetAdminTransactionHistory", status, from, to).Return(expectedData, nil).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetAdminTransactionHistory(status, from, to)
		assert.Nil(t, err)
		assert.Equal(t, len(res), len(expectedData))
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("GetAdminTransactionHistory", status, from, to).Return([]transaction.AdmTransactionRes{}, errors.New("server problem")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetAdminTransactionHistory(status, from, to)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, 0, len(res))
		data.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("GetAdminTransactionHistory", status, from, to).Return([]transaction.AdmTransactionRes{}, errors.New("data not found")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetAdminTransactionHistory(status, from, to)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, 0, len(res))
		data.AssertExpectations(t)
	})
}

func TestGetAdminTransactionDetails(t *testing.T) {
	data := mocks.NewTransactionData(t)
	transactionId := 1
	expectedData := transaction.AdmTransactionResDet{
		ID:                1,
		TenantId:          1,
		TenantName:        "customer1",
		TotalBill:         495000,
		CreatedAt:         ToTime("2023-01-26T02:11:48"),
		TransactionStatus: "success",
		InvoiceUrl:        "https://mediasosial.s3.ap-southeast-1.amazonaws.com/invoice/InvoiceSimple-PDF-Template.pdf",
		TransactionProductRes: []transaction.TransactionProductRes{
			{
				ProductId: 1,
				Quantity:  15,
				Price:     5000,
			},
			{
				ProductId: 2,
				Quantity:  10,
				Price:     20000,
			},
			{
				ProductId: 7,
				Quantity:  20,
				Price:     12000,
			},
		},
	}

	t.Run("success get admin transaction detail", func(t *testing.T) {
		data.On("GetAdminTransactionDetails", uint(transactionId)).Return(expectedData, nil).Once()
		srv := New(data)

		res, err := srv.GetAdminTransactionDetails(uint(transactionId))
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("GetAdminTransactionDetails", uint(transactionId)).Return(transaction.AdmTransactionResDet{}, errors.New("server problem")).Once()
		srv := New(data)

		res, err := srv.GetAdminTransactionDetails(uint(transactionId))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})
	t.Run("data not found", func(t *testing.T) {
		data.On("GetAdminTransactionDetails", uint(transactionId)).Return(transaction.AdmTransactionResDet{}, errors.New("data not found")).Once()
		srv := New(data)

		res, err := srv.GetAdminTransactionDetails(uint(transactionId))
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, res.TransactionStatus, "")
		data.AssertExpectations(t)
	})
}

func TestNotificationTransactionStatus(t *testing.T) {
	data := mocks.NewTransactionData(t)
	invNo := "INV-20230206-SELL-159"
	// var test *coreapi.TransactionStatusResponse
	t.Run("success handling payment notification", func(t *testing.T) {
		c := config.MidtransCoreAPIClient()

		// 4. Check transaction to Midtrans with param invoice number
		transactionStatusResp, _ := c.CheckTransaction(invNo)
		data.On("NotificationTransactionStatus", invNo, transactionStatusResp.TransactionStatus).Return(nil).Once()
		srv := New(data)

		err := srv.NotificationTransactionStatus(invNo)
		assert.Nil(t, err)
		data.AssertExpectations(t)
	})

	t.Run("error check transaction status", func(t *testing.T) {
		srv := New(data)

		err := srv.NotificationTransactionStatus("xxxxx")
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "error check")
	})

	// t.Run("error calling data", func(t *testing.T) {
	// 	repo.On("CheckTransaction", "Transaction-1").Return(test, nil)
	// 	repo.On("NotificationTransactionStatus", "xxxx", "failure").Return(errors.New("error update order status"))
	// 	srv := New(repo)

	// 	err := srv.NotificationTransactionStatus("xxxx")
	// 	assert.NotNil(t, err)
	// 	assert.Contains(t, err.Error(), "error calling")
	// })

	// t.Run("error calling data", func(t *testing.T) {
	// 	repo.On("CheckTransaction", "Transaction-1").Return(test, nil)
	// 	repo.On("NotificationTransactionStatus", "xxxx", "failure").Return(nil)
	// 	srv := New(repo)

	// 	err := srv.NotificationTransactionStatus("xxxx")
	// 	assert.Nil(t, err)
	// 	repo.AssertExpectations(t)
	// })
}
