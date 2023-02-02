package migration

import (
	product "sirloinapi/features/product/data"
	user "sirloinapi/features/user/data"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(user.User{})
	db.AutoMigrate(product.Product{})
	// db.AutoMigrate(cart.Cart{})
	// if !db.Migrator().HasColumn(&cart.CartItem{}, "Qty") {
	// 	db.Migrator().AddColumn(&cart.CartItem{}, "Qty")
	// }
	// if !db.Migrator().HasColumn(&cart.CartItem{}, "DeletedAt") {
	// 	db.Migrator().AddColumn(&cart.CartItem{}, "DeletedAt")
	// }
}
