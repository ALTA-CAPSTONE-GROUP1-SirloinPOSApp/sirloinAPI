package data

import (
	"sirloinapi/features/customer"
	user "sirloinapi/features/user/data"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	UserId      uint
	Email       string `gorm:"unique"`
	PhoneNumber string `gorm:"unique"`
	Name        string
	Address     string

	User user.User `gorm:"foreignkey:UserId;association_foreignkey:ID"`
}

func ToCore(data Customer) customer.Core {
	return customer.Core{
		ID:          data.ID,
		Email:       data.Email,
		Name:        data.Name,
		Address:     data.Address,
		PhoneNumber: data.PhoneNumber,
	}
}

func CoreToData(data customer.Core) Customer {
	return Customer{
		Model:       gorm.Model{ID: data.ID},
		Email:       data.Email,
		Name:        data.Name,
		Address:     data.Address,
		PhoneNumber: data.PhoneNumber,
	}
}
