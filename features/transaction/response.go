package transaction

import (
	"time"
)

type TransactionResponse struct {
	ID                uint      `json:"id"`
	CustomerId        uint      `json:"customer_id"`
	CustomerName      string    `json:"customer_name"`
	TotalPrice        float64   `json:"total_price"`
	Discount          float64   `json:"discount"`
	TotalBill         float64   `json:"total_bill"`
	CreatedAt         time.Time `json:"created_at"`
	TransactionStatus string    `json:"transaction_Status"`
	InvoiceNumber     string    `json:"invoice_number"`
	InvoiceUrl        string    `json:"invoice_url"`
	PaymentUrl        string    `json:"payment_url"`
}
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
	PdfUrl            string    `json:"pdf_url"`
	UserEmail         string    `json:"user_email"`
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
type AdmTransactionResponse struct {
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
func ToAdmResponse(data AdmTransactionRes) AdmTransactionResponse {
	return AdmTransactionResponse{
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

func ToAdmRespArr(data []AdmTransactionRes) []AdmTransactionResponse {
	arrRes := []AdmTransactionResponse{}
	for _, v := range data {
		tmp := ToAdmResponse(v)
		arrRes = append(arrRes, tmp)
	}
	return arrRes
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
	SellerName      string    `json:"seller_name"`
	SellerEmail     string    `json:"seller_email"`
	SellerPhone     string    `json:"seller_phone"`
	SellerAddress   string    `json:"seller_address"`
	CustomerName    string    `json:"customer_name"`
	CustomerEmail   string    `json:"customer_email"`
	CustomerPhone   string    `json:"customer_phone"`
	CustomerAddress string    `json:"customer_address"`
	SubTotal        float64   `json:"sub_total"`
	Discount        float64   `json:"discount"`
	DiscountAmount  float64   `json:"discount_amount"`
	TotalPrice      float64   `json:"total_price"`
}

type TransactionDetInv struct {
	InvoiceNumber   string    `json:"invoice_number"`
	TransactionDate time.Time `json:"created_at"`
	SellerName      string    `json:"seller_name"`
	SellerEmail     string    `json:"seller_email"`
	SellerPhone     string    `json:"seller_phone"`
	SellerAddress   string    `json:"seller_address"`
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

func InvToDetail(t TransactionInv) TransactionDetInv {
	return TransactionDetInv{
		InvoiceNumber:   t.InvoiceNumber,
		TransactionDate: t.TransactionDate,
		SellerName:      t.SellerName,
		SellerEmail:     t.SellerEmail,
		SellerPhone:     t.SellerPhone,
		SellerAddress:   t.SellerAddress,
		CustomerName:    t.CustomerName,
		CustomerEmail:   t.CustomerEmail,
		CustomerPhone:   t.CustomerPhone,
		CustomerAddress: t.CustomerAddress,
		SubTotal:        t.SubTotal,
		Discount:        t.Discount,
		DiscountAmount:  t.DiscountAmount,
		TotalPrice:      t.TotalPrice,
	}
}

func ToResponse(data Core) TransactionResponse {
	return TransactionResponse{
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
func ToResponseArr(data []Core) []TransactionResponse {
	arrRes := []TransactionResponse{}
	for _, v := range data {
		tmp := ToResponse(v)
		arrRes = append(arrRes, tmp)
	}
	return arrRes
}
