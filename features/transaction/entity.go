package transaction

import (
	"time"

	"github.com/labstack/echo/v4"
)

type Core struct {
	ID            uint      `json:"id"`
	UserId        uint      `json:"user_id"`
	CustomerId    uint      `json:"customer_id"`
	CustomerName  string    `json:"customer_name"`
	TotalPrice    float64   `json:"total_price"`
	Discount      float64   `json:"discount"`
	TotalBill     float64   `json:"total_bill"`
	CreatedAt     time.Time `json:"created_at"`
	OrderStatus   string    `json:"order_Status"`
	ProductStatus string    `json:"product_status"`
	InvoiceNumber string    `json:"invoice_number"`
	InvoiceUrl    string    `json:"invoice_url"`
	PaymentUrl    string    `json:"payment_url"`
}

type TransactionHandler interface {
	AddSell() echo.HandlerFunc
	// GetOrderHistory() echo.HandlerFunc
	// GetSellingHistory() echo.HandlerFunc
	// NotificationTransactionStatus() echo.HandlerFunc
	// UpdateStatus() echo.HandlerFunc
}

type TransactionService interface {
	AddSell(token interface{}, uCart Cart) (Core, error)
	// GetOrderHistory(token interface{}) ([]Core, error)
	// GetSellingHistory(token interface{}) ([]Core, error)
	// NotificationTransactionStatus(transactionId string) error
	// UpdateStatus(orderid uint, status string) error
}

type TransactionData interface {
	AddSell(userId uint, uCart Cart) (Core, error)
	// GetOrderHistory(userId uint) ([]Core, error)
	// GetSellingHistory(userId uint) ([]Core, error)
	// NotificationTransactionStatus(transactionId, transStatus string) error
	// UpdateStatus(orderid uint, status string) error
}
