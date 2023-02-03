package data

import (
	"log"
	"sirloinapi/features/customer"

	"gorm.io/gorm"
)

type customerQry struct {
	db *gorm.DB
}

func New(db *gorm.DB) customer.CustomerData {
	return &customerQry{
		db: db,
	}
}

func (cq *customerQry) Add(userId uint, newCustomer customer.Core) (customer.Core, error) {
	cnv := CoreToData(newCustomer)
	cnv.UserId = userId
	err := cq.db.Create(&cnv).Error
	if err != nil {
		log.Println("error create query: ", err.Error())
		return customer.Core{}, err
	}

	newCustomer.ID = cnv.ID

	return newCustomer, nil
}
