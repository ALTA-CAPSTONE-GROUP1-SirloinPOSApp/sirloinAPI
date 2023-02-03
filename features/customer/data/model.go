package data

import (
	user "sirloinapi/features/user/data"

	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	UserId  uint
	Email   string
	Phone   string
	Name    string
	Address string

	User user.User `gorm:"foreignkey:UserId;association_foreignkey:ID"`
}
