package services

import (
	"errors"
	"sirloinapi/features/customer"
	"sirloinapi/helper"
	"sirloinapi/mocks"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	data := mocks.NewCustomerData(t)
	newCustomer := customer.Core{
		Email:       "mfauzanptra@gmail.com",
		Name:        "Muhamad Fauzan Putra",
		PhoneNumber: "085659171799",
		Address:     "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}
	expectedData := customer.Core{
		Email:       "mfauzanptra@gmail.com",
		Name:        "Muhamad Fauzan Putra",
		PhoneNumber: "085659171799",
		Address:     "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}
	t.Run("success add new customer", func(t *testing.T) {
		data.On("Add", uint(2), newCustomer).Return(expectedData, nil).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, newCustomer)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		assert.Equal(t, expectedData.Name, res.Name)
		data.AssertExpectations(t)
	})
	t.Run("server problem", func(t *testing.T) {
		data.On("Add", uint(2), newCustomer).Return(customer.Core{}, errors.New("server problem")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("Add", uint(2), newCustomer).Return(customer.Core{}, errors.New("data not found")).Once()

		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token error")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("email already exist", func(t *testing.T) {
		data.On("Add", uint(2), newCustomer).Return(customer.Core{}, errors.New("Duplicate customers.email")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user already exist")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("phone number already exist", func(t *testing.T) {
		data.On("Add", uint(2), newCustomer).Return(customer.Core{}, errors.New("Duplicate customers.phone_number")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "phone number already exist")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("Name validasi", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		newCustomer.Name = "AMR#%^%"
		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "alpha_space")
		assert.Equal(t, res.Name, "")
	})

	t.Run("Email validasi", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		newCustomer.Name = "Toko Ira"
		newCustomer.Email = "tokoira.com"
		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "email")
		assert.Equal(t, res.Name, "")
	})

	t.Run("Phone number validasi", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		newCustomer.Name = "Toko Ira"
		newCustomer.Email = "tokoira@gmail.com"
		newCustomer.PhoneNumber = "oiajsdojhasodi"
		res, err := srv.Add(pToken, newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "PhoneNumber")
		assert.Equal(t, res.Name, "")
	})

}

func TestUpdate(t *testing.T) {

	data := mocks.NewCustomerData(t)
	newCustomer := customer.Core{
		Email:       "mfauzanptra@gmail.com",
		Name:        "Muhamad Fauzan Putra",
		PhoneNumber: "085659171799",
		Address:     "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}

	expectedData := customer.Core{
		Email:       "mfauzanptra@gmail.com",
		Name:        "Muhamad Fauzan Putra",
		PhoneNumber: "085659171799",
		Address:     "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}

	userId := 2
	customerId := 1

	t.Run("success edit customer data", func(t *testing.T) {
		data.On("Update", uint(userId), uint(customerId), newCustomer).Return(expectedData, nil).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.Nil(t, err)
		assert.Equal(t, expectedData.ID, res.ID)
		assert.Equal(t, expectedData.Name, res.Name)
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("Update", uint(userId), uint(customerId), newCustomer).Return(customer.Core{}, errors.New("server problem")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("Update", uint(userId), uint(customerId), newCustomer).Return(customer.Core{}, errors.New("data not found")).Once()

		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token error")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("email already exist", func(t *testing.T) {
		data.On("Update", uint(userId), uint(customerId), newCustomer).Return(customer.Core{}, errors.New("Duplicate customers.email")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "user already exist")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("phone number already exist", func(t *testing.T) {
		data.On("Update", uint(userId), uint(customerId), newCustomer).Return(customer.Core{}, errors.New("Duplicate customers.phone_number")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "phone number already exist")
		assert.Equal(t, res.Name, "")
		data.AssertExpectations(t)
	})

	t.Run("Name validasi", func(t *testing.T) {
		srv := New(data)

		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		newCustomer.Name = "AMR#%^%"
		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "alpha_space")
		assert.Equal(t, res.Name, "")
	})

	t.Run("Email validasi", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		newCustomer.Name = "Toko Ira"
		newCustomer.Email = "tokoira.com"
		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "email")
		assert.Equal(t, res.Name, "")
	})

	t.Run("Phone number validasi", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		newCustomer.Name = "Toko Ira"
		newCustomer.Email = "tokoira@gmail.com"
		newCustomer.PhoneNumber = "oiajsdojhasodi"
		res, err := srv.Update(pToken, uint(customerId), newCustomer)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "PhoneNumber")
		assert.Equal(t, res.Name, "")
	})
}

func TestGetUserCustomers(t *testing.T) {

	data := mocks.NewCustomerData(t)
	expectedData := []customer.Core{{
		Email:       "mfauzanptra@gmail.com",
		Name:        "Muhamad Fauzan Putra",
		PhoneNumber: "085659171799",
		Address:     "Jln. Lembayung No 24, Bantul, Yogyakarta",
	}, {
		Email:       "arimrizal94@gmail.com",
		Name:        "Ari Muhammad Rizal",
		PhoneNumber: "081230255973",
		Address:     "Kab. Kediri, Jawa Timur",
	}}
	userId := 2

	t.Run("success get all customers", func(t *testing.T) {
		data.On("GetUserCustomers", uint(userId)).Return(expectedData, nil).Once()
		srv := New(data)
		_, token := helper.GenerateJWT(2)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetUserCustomers(pToken)
		assert.Nil(t, err)
		assert.Equal(t, len(res), len(expectedData))
		data.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		data.On("GetUserCustomers", uint(userId)).Return([]customer.Core{}, errors.New("server problem")).Once()
		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetUserCustomers(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, 0, len(res))
		data.AssertExpectations(t)
	})

	t.Run("data not found", func(t *testing.T) {
		data.On("GetUserCustomers", uint(userId)).Return([]customer.Core{}, errors.New("data not found")).Once()

		srv := New(data)

		_, token := helper.GenerateJWT(userId)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetUserCustomers(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, 0, len(res))
		data.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(data)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.GetUserCustomers(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "token error")
		assert.Equal(t, 0, len(res))
	})
}
