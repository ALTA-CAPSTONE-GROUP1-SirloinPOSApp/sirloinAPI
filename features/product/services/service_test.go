package services

import (
	"errors"
	"mime/multipart"
	"sirloinapi/features/product"
	"sirloinapi/helper"
	"sirloinapi/mocks"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdd(t *testing.T) {
	repo := mocks.NewProductData(t)

	inputData := product.Core{
		ProductName:  "Indomie goreng",
		Upc:          "85191919891918519",
		Category:     "makanan",
		MinimumStock: 5,
		BuyingPrice:  2000,
		Price:        3000,
		Supplier:     "saya",
		Stock:        20,
		// ProductImage: "Indomie.jpg",
	}
	inputDataVld := product.Core{
		ProductName:  "Indomie goreng",
		Upc:          "85191919891918519",
		Category:     "makanan",
		MinimumStock: 5,
		BuyingPrice:  2000,
		Price:        3000,
		Supplier:     "saya",
		// ProductImage: "Indomie.jpg",
	}
	resData := product.Core{
		ProductName:  "Indomie goreng",
		Upc:          "85191919891918519",
		Category:     "makanan",
		MinimumStock: 5,
		BuyingPrice:  2000,
		Price:        3000,
		Supplier:     "saya",
		Stock:        20,
		// ProductImage: "https://socmedapibucket.s3.ap-southeast-1.amazonaws.com/files/post/1/indomie-photo.jpeg",
	}
	var a *multipart.FileHeader

	t.Run("success add post", func(t *testing.T) {

		repo.On("Add", uint(1), inputData, mock.Anything).Return(resData, nil).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Add(pToken, inputData, a)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("user not found ", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		res, err := srv.Add(token, inputData, a)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data not found", func(t *testing.T) {
		repo.On("Add", uint(1), inputData, a).Return(product.Core{}, errors.New("data not found")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, inputData, a)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		repo.On("Add", uint(1), inputData, a).Return(product.Core{}, errors.New("server problem"))
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, inputData, a)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Add(pToken, inputDataVld, a)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
	})
}

func TestUpdate(t *testing.T) {
	repo := mocks.NewProductData(t)

	inputData := product.Core{
		ProductName:  "Indomie goreng",
		Upc:          "85191919891918519",
		Category:     "makanan",
		MinimumStock: 5,
		BuyingPrice:  2000,
		Price:        3000,
		Supplier:     "saya",
		Stock:        20,
		// ProductImage: "Indomie.jpg",
	}
	resData := product.Core{
		ProductName:  "Indomie goreng",
		Upc:          "85191919891918519",
		Category:     "makanan",
		MinimumStock: 5,
		BuyingPrice:  2000,
		Price:        3000,
		Supplier:     "saya",
		Stock:        20,
		// ProductImage: "https://socmedapibucket.s3.ap-southeast-1.amazonaws.com/files/post/1/indomie-photo.jpeg",
	}
	var a *multipart.FileHeader
	productId := uint(1)

	t.Run("success update product", func(t *testing.T) {
		repo.On("Update", uint(1), productId, inputData, mock.Anything).Return(resData, nil).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.Update(pToken, productId, inputData, a)
		assert.Nil(t, err)
		assert.Equal(t, resData.ID, res.ID)
		repo.AssertExpectations(t)
	})

	t.Run("user not found ", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		res, err := srv.Update(token, productId, inputData, a)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, uint(0), res.ID)
	})

	t.Run("data not found", func(t *testing.T) {
		repo.On("Update", uint(1), productId, inputData, a).Return(product.Core{}, errors.New("data not found")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, productId, inputData, a)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		repo.On("Update", uint(1), productId, inputData, a).Return(product.Core{}, errors.New("server problem"))
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		res, err := srv.Update(pToken, productId, inputData, a)
		assert.NotNil(t, err)
		assert.Equal(t, uint(0), res.ID)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	productId := uint(1)
	repo := mocks.NewProductData(t)

	t.Run("success delete product", func(t *testing.T) {
		repo.On("Delete", uint(1), productId).Return(nil).Once()
		srv := New(repo)
		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		err := srv.Delete(pToken, productId)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("user not found ", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		err := srv.Delete(token, productId)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("data not found", func(t *testing.T) {
		repo.On("Delete", uint(1), productId).Return(errors.New("data not found")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, productId)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		repo.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		repo.On("Delete", uint(1), productId).Return(errors.New("server problem")).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true
		err := srv.Delete(pToken, productId)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		repo.AssertExpectations(t)
	})
}

func TestGetAllProducts(t *testing.T) {
	repo := mocks.NewProductData(t)
	resdata := []product.Core{
		{
			ID:           1,
			UserId:       1,
			UserName:     "Fauzan",
			ProductName:  "Adidas NMD",
			ProductImage: "url",
			Stock:        10,
			Price:        900000,
			Upc:          "15191981981981",
			Category:     "food",
		},
	}

	t.Run("success get all user products", func(t *testing.T) {
		repo.On("GetUserProducts", uint(1)).Return(resdata, nil).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetUserProducts(pToken)
		assert.Nil(t, err)
		assert.Equal(t, len(res), len(resdata))
		repo.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(0)

		_, err := srv.GetUserProducts(token)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("data not found", func(t *testing.T) {
		repo.On("GetUserProducts", uint(1)).Return([]product.Core{}, errors.New("data not found")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		srv := New(repo)

		res, err := srv.GetUserProducts(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, 0, len(res))
		repo.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		repo.On("GetUserProducts", uint(1)).Return([]product.Core{}, errors.New("server problem")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		srv := New(repo)

		res, err := srv.GetUserProducts(pToken)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, 0, len(res))
		repo.AssertExpectations(t)
	})
}

func TestGetProductById(t *testing.T) {
	repo := mocks.NewProductData(t)
	productId := uint(1)

	resdata := product.Core{
		ID:           1,
		UserId:       1,
		UserName:     "Fauzan",
		ProductName:  "Adidas NMD",
		ProductImage: "url",
		Stock:        10,
		Price:        900000,
		Upc:          "15191981981981",
		Category:     "food",
	}

	t.Run("success get product detail by id", func(t *testing.T) {
		repo.On("GetProductById", uint(1), productId).Return(resdata, nil).Once()
		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetProductById(pToken, productId)
		assert.Nil(t, err)
		assert.Equal(t, res.UserName, resdata.UserName)
		repo.AssertExpectations(t)
	})

	t.Run("jwt not valid", func(t *testing.T) {
		srv := New(repo)

		_, token := helper.GenerateJWT(1)

		_, err := srv.GetProductById(token, productId)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
	})

	t.Run("data not found", func(t *testing.T) {
		repo.On("GetProductById", uint(1), productId).Return(product.Core{}, errors.New("data not found")).Once()

		srv := New(repo)

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		res, err := srv.GetProductById(pToken, productId)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "not found")
		assert.Equal(t, res.ID, uint(0))
		repo.AssertExpectations(t)
	})

	t.Run("server problem", func(t *testing.T) {
		repo.On("GetProductById", uint(1), productId).Return(product.Core{}, errors.New("server problem")).Once()

		_, token := helper.GenerateJWT(1)
		pToken := token.(*jwt.Token)
		pToken.Valid = true

		srv := New(repo)

		res, err := srv.GetProductById(pToken, productId)
		assert.NotNil(t, err)
		assert.ErrorContains(t, err, "server")
		assert.Equal(t, res.ID, uint(0))
		repo.AssertExpectations(t)
	})
}
