package data

import (
	"sirloinapi/features/user"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	BusinessName string
	Email        string `gorm:"unique"`
	Address      string
	PhoneNumber  string `gorm:"unique"`
	Password     string
}

type DeviceToken struct {
	UserId uint
	Token  string
	User   User `gorm:"foreignkey:UserId;association_foreignkey:ID"`
}

func ToCore(data User) user.Core {
	return user.Core{
		ID:           data.ID,
		Email:        data.Email,
		BusinessName: data.BusinessName,
		Address:      data.Address,
		PhoneNumber:  data.PhoneNumber,
		Password:     data.Password,
	}
}

func CoreToData(data user.Core) User {
	return User{
		Model:        gorm.Model{ID: data.ID},
		Email:        data.Email,
		BusinessName: data.BusinessName,
		Address:      data.Address,
		PhoneNumber:  data.PhoneNumber,
		Password:     data.Password,
	}
}
