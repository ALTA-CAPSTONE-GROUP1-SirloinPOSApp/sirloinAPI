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

type AdmTransactionRes struct {
	ID                uint      `json:"id"`
	TenantId          uint      `json:"tenant_id"`
	TenantName        string    `json:"tenant_name"`
	TotalBill         float64   `json:"total_bill"`
	CreatedAt         time.Time `json:"created_at"`
	TransactionStatus string    `json:"transaction_status"`
	InvoiceNumber     string    `json:"invoice_number"`
	InvoiceUrl        string    `json:"invoice_url"`
	PaymentUrl        string    `json:"payment_url"`
}

type AdmTransactionResDet struct {
	ID                    uint      `json:"id"`
	TenantId              uint      `json:"tenant_id"`
	TenantName            string    `json:"tenant_name"`
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

func ToAdmResp(data AdmTransactionRes) AdmTransactionResDet {
	return AdmTransactionResDet{
		ID:                data.ID,
		TenantId:          data.TenantId,
		TenantName:        data.TenantName,
		TotalBill:         data.TotalBill,
		CreatedAt:         data.CreatedAt,
		TransactionStatus: data.TransactionStatus,
		InvoiceNumber:     data.InvoiceNumber,
		InvoiceUrl:        data.InvoiceUrl,
		PaymentUrl:        data.PaymentUrl,
	}
}

type ItemsInv struct {
	ItemName   string  `json:"item_name"`
	Quantity   int     `json:"quantity"`
	Price      float64 `json:"price"`
	TotalPrice float64 `json:"total_price"`
}

type TransactionInv struct {
	InvoiceNumber   string    `json:"invoice_number"`
	TransactionDate time.Time `json:"created_at"`
	SellerName      string    `json:"tenant_name"`
	SellerPhone     string    `json:"tenant_phone"`
	SellerAddress   string    `json:"tenant_address"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerPhone   string    `json:"customer_phone"`
	CustomerAddress string    `json:"customer_address"`
	SubTotal        float64   `json:"sub_total"`
	Discount        float64   `json:"discount"`
	DiscountAmount  float64   `json:"discount_amount"`
	TotalPrice      float64   `json:"total_price"`
	Items           []ItemsInv
}
