package handler

import "sirloinapi/features/user"

type RegisterRequest struct {
	BusinessName string `json:"business_name" form:"business_name"`
	Email        string `json:"email" form:"email"`
	Address      string `json:"address" form:"address"`
	PhoneNumber  string `json:"phone_number" form:"phone_number"`
	Password     string `json:"password" form:"password"`
}

type LoginReqest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

func ToCore(data interface{}) *user.Core {
	res := user.Core{}

	switch docs := data.(type) {
	case RegisterRequest:
		cnv := docs
		res.Email = cnv.Email
		res.BusinessName = cnv.BusinessName
		res.PhoneNumber = cnv.PhoneNumber
		res.Password = cnv.Password
		res.Address = cnv.Address
	case LoginReqest:
		cnv := docs
		res.Email = cnv.Email
		res.Password = cnv.Password
	default:
		return nil
	}
	return &res
}
