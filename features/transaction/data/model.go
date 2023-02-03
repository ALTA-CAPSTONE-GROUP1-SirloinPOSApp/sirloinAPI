package data

import (
	customer "sirloinapi/features/customer/data"
	product "sirloinapi/features/product/data"
	"sirloinapi/features/transaction"
	user "sirloinapi/features/user/data"
	"time"
)

type Transaction struct {
	ID                uint `gorm:"PRIMARY_KEY;AUTO_INCREMENT;NOT NULL"`
	UserId            uint
	CustomerId        uint
	TotalPrice        float64
	Discount          float64
	TotalBill         float64
	CreatedAt         time.Time
	TransactionStatus string
	ProductStatus     string
	InvoiceNumber     string
	InvoiceUrl        string
	PaymentUrl        string

	Customer customer.Customer `gorm:"foreignkey:CustomerId;association_foreignkey:ID"`
	User     user.User         `gorm:"foreignkey:UserId;association_foreignkey:ID"`
}

type TransactionProduct struct {
	TransactionId uint
	ProductId     uint
	Quantity      int
	Price         float64
	TotalPrice    float64

	Transaction Transaction     `gorm:"foreignkey:TransactionId;association_foreignkey:ID"`
	Product     product.Product `gorm:"foreignkey:ProductId;association_foreignkey:ID"`
}

func DataToCoreT(data Transaction) transaction.Core {
	return transaction.Core{
		ID:                data.ID,
		UserId:            data.UserId,
		CustomerId:        data.CustomerId,
		TotalPrice:        data.TotalPrice,
		Discount:          data.Discount,
		TotalBill:         data.TotalBill,
		CreatedAt:         data.CreatedAt,
		TransactionStatus: data.TransactionStatus,
		ProductStatus:     data.ProductStatus,
		InvoiceNumber:     data.InvoiceNumber,
		InvoiceUrl:        data.InvoiceUrl,
		PaymentUrl:        data.PaymentUrl,
	}
}
