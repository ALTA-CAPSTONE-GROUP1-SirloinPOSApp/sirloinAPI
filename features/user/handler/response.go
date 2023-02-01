package handler

import "sirloinapi/features/user"

type UserResponse struct {
	ID           uint   `json:"id" form:"id"`
	BusinessName string `json:"business_name" form:"business_name"`
	Email        string `json:"email" form:"email"`
	Address      string `json:"address" form:"address"`
	PhoneNumber  string `json:"phone_number" form:"phone_number"`
}

type LoginResp struct {
	ID           uint   `json:"id" form:"id"`
	BusinessName string `json:"business_name" form:"business_name"`
	Email        string `json:"email" form:"email"`
	Address      string `json:"address" form:"address"`
	PhoneNumber  string `json:"phone_number" form:"phone_number"`
	Token        string `json:"token" form:"token"`
}

func ToResponse(data user.Core) UserResponse {
	return UserResponse{
		ID:           data.ID,
		Email:        data.Email,
		BusinessName: data.BusinessName,
		PhoneNumber:  data.PhoneNumber,
		Address:      data.Address,
	}
}

func ToLoginResp(data user.Core, token string) LoginResp {
	return LoginResp{
		ID:           data.ID,
		Email:        data.Email,
		BusinessName: data.BusinessName,
		PhoneNumber:  data.PhoneNumber,
		Address:      data.Address,
		Token:        token,
	}
}
