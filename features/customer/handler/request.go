package handler

import (
	"sirloinapi/features/customer"
)

type AddCustomerRequest struct {
	Name        string `json:"name" form:"name"`
	Email       string `json:"email" form:"email"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

func ToCore(data interface{}) *customer.Core {
	res := customer.Core{}

	switch docs := data.(type) {
	case AddCustomerRequest:
		cnv := docs
		res.Email = cnv.Email
		res.Name = cnv.Name
		res.PhoneNumber = cnv.PhoneNumber
		res.Address = cnv.Address
	default:
		return nil
	}
	return &res
}
