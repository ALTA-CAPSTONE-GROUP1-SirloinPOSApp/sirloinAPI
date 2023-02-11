package data

import (
	"sirloinapi/features/customer"
	user "sirloinapi/features/user/data"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	UserId      uint
	Email       string
	PhoneNumber string
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

func ToCoreArr(data []Customer) []customer.Core {
	arrRes := []customer.Core{}
	for _, v := range data {
		tmp := ToCore(v)
		arrRes = append(arrRes, tmp)
	}
	return arrRes
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
