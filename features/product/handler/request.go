package handler

import "sirloinapi/features/product"

type AddProductReq struct {
	Upc          string  `json:"upc" form:"upc"`
	Category     string  `json:"category" form:"category"`
	ProductName  string  `json:"product_name" form:"product_name"`
	Stock        int     `json:"stock" form:"stock"`
	MinimumStock int     `json:"minimum_stock" form:"minimum_stock"`
	BuyingPrice  float64 `json:"buying_price" form:"buying_price"`
	Price        float64 `json:"price" form:"price"`
	ProductImage string  `json:"product_image" form:"product_image"`
	Supplier     string  `json:"supplier" form:"supplier"`
}

func ToCore(data interface{}) *product.Core {
	res := product.Core{}

	switch docs := data.(type) {
	case AddProductReq:
		cnv := docs
		res.Upc = cnv.Upc
		res.Category = cnv.Category
		res.ProductName = cnv.ProductName
		res.Stock = cnv.Stock
		res.MinimumStock = cnv.MinimumStock
		res.BuyingPrice = cnv.BuyingPrice
		res.Price = cnv.Price
		res.ProductImage = cnv.ProductImage
		res.Supplier = cnv.Supplier
	default:
		return nil
	}
	return &res
}
