package config

import (
	"fmt"
	"log"
	user "sirloinapi/features/user/data"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB(ac AppConfig) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		ac.DBUser, ac.DBPass, ac.DBHost, ac.DBPort, ac.DBName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("database connection error : ", err.Error())
		return nil
	}

	return db
}

func Migrate(db *gorm.DB) {
	db.AutoMigrate(user.User{})
	// db.AutoMigrate(item.Item{})
	// db.AutoMigrate(cart.Cart{})
	// if !db.Migrator().HasColumn(&cart.CartItem{}, "Qty") {
	// 	db.Migrator().AddColumn(&cart.CartItem{}, "Qty")
	// }
	// if !db.Migrator().HasColumn(&cart.CartItem{}, "DeletedAt") {
	// 	db.Migrator().AddColumn(&cart.CartItem{}, "DeletedAt")
	// }
}
