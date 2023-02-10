package data

import (
	"sirloinapi/features/product"
	user "sirloinapi/features/user/data"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	UserId       uint
	User         user.User
	Upc          string `gorm:"unique"`
	Category     string
	ProductName  string `gorm:"unique"`
	Stock        int
	MinimumStock int
	BuyingPrice  float64
	Price        float64
	ProductImage string
	Supplier     string
	ItemsSold    int
}

func DataToCore(data Product) product.Core {
	return product.Core{
		ID:           data.ID,
		UserId:       data.UserId,
		Upc:          data.Upc,
		Category:     data.Category,
		ProductName:  data.ProductName,
		Stock:        data.Stock,
		MinimumStock: data.MinimumStock,
		BuyingPrice:  data.BuyingPrice,
		Price:        data.Price,
		ProductImage: data.ProductImage,
		Supplier:     data.Supplier,
		ItemsSold:    data.ItemsSold,
	}
}

func CoreToData(data product.Core) Product {
	return Product{
		Model:        gorm.Model{ID: data.ID},
		UserId:       data.UserId,
		Upc:          data.Upc,
		Category:     data.Category,
		ProductName:  data.ProductName,
		Stock:        data.Stock,
		MinimumStock: data.MinimumStock,
		BuyingPrice:  data.BuyingPrice,
		Price:        data.Price,
		ProductImage: data.ProductImage,
		Supplier:     data.Supplier,
		ItemsSold:    data.ItemsSold,
	}
}
