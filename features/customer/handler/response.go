package handler

type CustomerResponse struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"business_name" form:"business_name"`
	Email       string `json:"email" form:"email"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}
