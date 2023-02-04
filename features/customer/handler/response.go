package handler

import "sirloinapi/features/customer"

type CustomerResponse struct {
	ID          uint   `json:"id" form:"id"`
	Name        string `json:"business_name" form:"business_name"`
	Email       string `json:"email" form:"email"`
	Address     string `json:"address" form:"address"`
	PhoneNumber string `json:"phone_number" form:"phone_number"`
}

func ToResponse(data customer.Core) CustomerResponse {
	return CustomerResponse{
		ID:          data.ID,
		Name:        data.Name,
		Email:       data.Email,
		Address:     data.Address,
		PhoneNumber: data.PhoneNumber,
	}
}

func ToResponseArr(data []customer.Core) []CustomerResponse {
	res := []CustomerResponse{}
	for _, v := range data {
		tmp := ToResponse(v)
		res = append(res, tmp)
	}
	return res
}
