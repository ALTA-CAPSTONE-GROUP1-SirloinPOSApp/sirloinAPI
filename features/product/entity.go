package product

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID           uint    `json:"id" form:"id"`
	UserId       uint    `json:"user_id" form:"user_id"`
	UserName     string  `json:"user_name" form:"user_name"`
	Upc          string  `json:"upc" form:"upc"`
	Category     string  `json:"category" form:"category"`
	ProductName  string  `json:"product_name" form:"product_name"`
	Stock        int     `json:"stock" form:"stock"`
	MinimumStock int     `json:"minimum_stock" form:"minimum_stock"`
	BuyingPrice  float64 `json:"buying_price" form:"buying_price"`
	Price        float64 `json:"price" form:"price"`
	ProductImage string  `json:"product_image" form:"product_image"`
	Supplier     string  `json:"supplier" form:"supplier"`
	ItemsSold    int     `json:"items_sold" form:"items_sold"`
}

type ProductHandler interface {
	Add() echo.HandlerFunc
	Update() echo.HandlerFunc
	Delete() echo.HandlerFunc
	GetUserProducts() echo.HandlerFunc
	GetProductById() echo.HandlerFunc
}

type ProductService interface {
	Add(token interface{}, newProduct Core, productImage *multipart.FileHeader) (Core, error)
	Update(token interface{}, productId uint, updProduct Core, productImage *multipart.FileHeader) (Core, error)
	Delete(token interface{}, productId uint) error
	GetUserProducts(token interface{}) ([]Core, error)
	GetProductById(token interface{}, productId uint) (Core, error)
}

type ProductData interface {
	Add(userId uint, newProduct Core, productImage *multipart.FileHeader) (Core, error)
	Update(userId, productId uint, updProduct Core, productImage *multipart.FileHeader) (Core, error)
	Delete(userId, productId uint) error
	GetUserProducts(userId uint) ([]Core, error)
	GetProductById(userid, productId uint) (Core, error)
}
