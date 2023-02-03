package services

import (
	"errors"
	"sirloinapi/features/transaction"
	"sirloinapi/helper"
	"sirloinapi/mocks"
	"testing"

	"github.com/golang-jwt/jwt"
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
