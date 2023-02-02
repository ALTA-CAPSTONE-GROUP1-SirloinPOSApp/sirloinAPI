package services

import (
	"errors"
	"log"
	"mime/multipart"
	"sirloinapi/features/product"
	"sirloinapi/helper"
	"strings"

	"github.com/go-playground/validator/v10"
)

type productSvc struct {
	qry product.ProductData
	vld *validator.Validate
}

func New(data product.ProductData) product.ProductService {
	return &productSvc{
		qry: data,
		vld: validator.New(),
	}
}

func (ps *productSvc) Add(token interface{}, newProduct product.Core, productImage *multipart.FileHeader) (product.Core, error) {
	err := ps.vld.Struct(newProduct)
	if err != nil {
		return product.Core{}, err
	}
	userId := helper.ExtractToken(token)
	if userId <= 0 {
		log.Println("\t error extract token add product")
		return product.Core{}, errors.New("user not found")
	}

	res, err := ps.qry.Add(uint(userId), newProduct, productImage)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		log.Println("\terror add query in service: ", err.Error())
		return product.Core{}, errors.New(msg)
	}

	return res, nil
}
func (ps *productSvc) Update(token interface{}, productId uint, updProduct product.Core, productImage *multipart.FileHeader) (product.Core, error) {
	userId := helper.ExtractToken(token)
	if userId <= 0 {
		log.Println("\t error extract token add product")
		return product.Core{}, errors.New("user not found")
	}

	res, err := ps.qry.Update(uint(userId), productId, updProduct, productImage)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "product data not found"
		} else {
			msg = "server problem"
		}
		log.Println("\terror update data in service: ", err.Error())
		return product.Core{}, errors.New(msg)
	}

	return res, nil
}
func (ps *productSvc) Delete(token interface{}, productId uint) error {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		log.Println("\terror extract token delete product service")
		return errors.New("user not found")
	}

	err := ps.qry.Delete(uint(userID), productId)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "product not found"
		} else {
			msg = "server problem"
		}
		log.Println("\terror calling delete data in service: ", err.Error())
		return errors.New(msg)
	}
	return nil
}
func (ps *productSvc) GetUserProducts(token interface{}) ([]product.Core, error) {
	userID := helper.ExtractToken(token)
	if userID <= 0 {
		log.Println("\terror extract token getUserProduct product service")
		return []product.Core{}, errors.New("user not found")
	}

	res, err := ps.qry.GetUserProducts(uint(userID))
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		return []product.Core{}, errors.New(msg)
	}
	return res, nil
}
func (ps *productSvc) GetProductById(token interface{}, productId uint) (product.Core, error) {
	userId := helper.ExtractToken(token)
	if userId <= 0 {
		log.Println("\terror extract token delete product service")
		return product.Core{}, errors.New("user not found")
	}

	res, err := ps.qry.GetProductById(uint(userId), productId)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		return product.Core{}, errors.New(msg)
	}
	return res, nil
}
