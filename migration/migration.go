package migration

import (
	customer "sirloinapi/features/customer/data"
	product "sirloinapi/features/product/data"
	trans "sirloinapi/features/transaction/data"
	user "sirloinapi/features/user/data"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(user.User{})
	db.AutoMigrate(user.DeviceToken{})
	db.AutoMigrate(product.Product{})
	db.AutoMigrate(customer.Customer{})
	db.AutoMigrate(trans.Transaction{})
	db.AutoMigrate(trans.TransactionProduct{})
}
