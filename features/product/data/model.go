package data

import (
	"sirloinapi/features/product"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	UserId       uint    `json:"user_id" form:"user_id"`
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
