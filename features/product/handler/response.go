package handler

import (
	"sirloinapi/features/product"
)

type GetProdResp struct {
	Id           uint    `json:"id"`
	Upc          string  `json:"upc"`
	Category     string  `json:"category"`
	ProductName  string  `json:"product_name"`
	Stock        int     `json:"stock"`
	MinimumStock int     `json:"minimum_stock"`
	BuyingPrice  float64 `json:"buying_price"`
	Price        float64 `json:"price"`
	ProductImage string  `json:"product_image"`
	Supplier     string  `json:"supplier"`
	ItemsSold    int     `json:"items_sold"`
}

func ToGetProdResp(p product.Core) GetProdResp {
	return GetProdResp{
		Id:           p.ID,
		Upc:          p.Upc,
		Category:     p.Category,
		ProductName:  p.ProductName,
		Stock:        p.Stock,
		MinimumStock: p.MinimumStock,
		BuyingPrice:  p.BuyingPrice,
		Price:        p.Price,
		ProductImage: p.ProductImage,
		Supplier:     p.Supplier,
		ItemsSold:    p.ItemsSold,
	}
}

func ToGetProdsResp(p []product.Core) []GetProdResp {
	res := []GetProdResp{}

	for _, v := range p {
		tmp := ToGetProdResp(v)
		res = append(res, tmp)
	}

	return res
}
