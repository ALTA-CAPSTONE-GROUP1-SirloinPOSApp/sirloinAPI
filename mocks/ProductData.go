// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	multipart "mime/multipart"

	mock "github.com/stretchr/testify/mock"

	product "sirloinapi/features/product"
)

// ProductData is an autogenerated mock type for the ProductData type
type ProductData struct {
	mock.Mock
}

// Add provides a mock function with given fields: userId, newProduct, productImage
func (_m *ProductData) Add(userId uint, newProduct product.Core, productImage *multipart.FileHeader) (product.Core, error) {
	ret := _m.Called(userId, newProduct, productImage)

	var r0 product.Core
	if rf, ok := ret.Get(0).(func(uint, product.Core, *multipart.FileHeader) product.Core); ok {
		r0 = rf(userId, newProduct, productImage)
	} else {
		r0 = ret.Get(0).(product.Core)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, product.Core, *multipart.FileHeader) error); ok {
		r1 = rf(userId, newProduct, productImage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: userId, productId
func (_m *ProductData) Delete(userId uint, productId uint) error {
	ret := _m.Called(userId, productId)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint, uint) error); ok {
		r0 = rf(userId, productId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetProductById provides a mock function with given fields: productId
func (_m *ProductData) GetProductById(productId uint) (product.Core, error) {
	ret := _m.Called(productId)

	var r0 product.Core
	if rf, ok := ret.Get(0).(func(uint) product.Core); ok {
		r0 = rf(productId)
	} else {
		r0 = ret.Get(0).(product.Core)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(productId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserProducts provides a mock function with given fields: userId
func (_m *ProductData) GetUserProducts(userId uint) ([]product.Core, error) {
	ret := _m.Called(userId)

	var r0 []product.Core
	if rf, ok := ret.Get(0).(func(uint) []product.Core); ok {
		r0 = rf(userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]product.Core)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: userId, productId, updProduct, productImage
func (_m *ProductData) Update(userId uint, productId uint, updProduct product.Core, productImage *multipart.FileHeader) (product.Core, error) {
	ret := _m.Called(userId, productId, updProduct, productImage)

	var r0 product.Core
	if rf, ok := ret.Get(0).(func(uint, uint, product.Core, *multipart.FileHeader) product.Core); ok {
		r0 = rf(userId, productId, updProduct, productImage)
	} else {
		r0 = ret.Get(0).(product.Core)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint, uint, product.Core, *multipart.FileHeader) error); ok {
		r1 = rf(userId, productId, updProduct, productImage)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewProductData interface {
	mock.TestingT
	Cleanup(func())
}

// NewProductData creates a new instance of ProductData. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewProductData(t mockConstructorTestingTNewProductData) *ProductData {
	mock := &ProductData{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
