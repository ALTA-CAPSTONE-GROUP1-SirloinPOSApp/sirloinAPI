package transaction

type Item struct {
	ProductId uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
type Cart struct {
	Items         []Item `json:"items"`
	CustomerId    uint   `json:"customer_id"`
	PaymentMethod string `json:"payment_method"`
}
