package transaction

import (
	"time"
)

type TransactionRes struct {
	ID                    uint      `json:"id"`
	CustomerId            uint      `json:"customer_id"`
	CustomerName          string    `json:"customer_name"`
	TotalPrice            float64   `json:"total_price"`
	Discount              float64   `json:"discount"`
	TotalBill             float64   `json:"total_bill"`
	CreatedAt             time.Time `json:"created_at"`
	TransactionStatus     string    `json:"transaction_status"`
	InvoiceNumber         string    `json:"invoice_number"`
	InvoiceUrl            string    `json:"invoice_url"`
	PaymentUrl            string    `json:"payment_url"`
	TransactionProductRes []TransactionProductRes
}

type TransactionProductRes struct {
	ProductId     uint    `json:"product_id"`
	ProductName   string  `json:"product_name"`
	Price         float64 `json:"price"`
	Quantity      int     `json:"quantity"`
	TotalPrice    float64 `json:"total_price"`
	Product_image string  `json:"product_image"`
}

func CoreToResp(data Core) TransactionRes {
	return TransactionRes{
		ID:                data.ID,
		CustomerId:        data.CustomerId,
		CustomerName:      data.CustomerName,
		TotalPrice:        data.TotalPrice,
		Discount:          data.Discount,
		TotalBill:         data.TotalBill,
		CreatedAt:         data.CreatedAt,
		TransactionStatus: data.TransactionStatus,
		InvoiceNumber:     data.InvoiceNumber,
		InvoiceUrl:        data.InvoiceUrl,
		PaymentUrl:        data.PaymentUrl,
	}
}
