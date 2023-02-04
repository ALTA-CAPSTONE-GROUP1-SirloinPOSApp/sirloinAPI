package data

import (
	"errors"
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

func (cq *customerQry) Update(userId, customerId uint, updateData customer.Core) (customer.Core, error) {
	cnvUpd := CoreToData(updateData)
	cnvUpd.UserId = userId
	qry := cq.db.Where("id = ? AND user_id = ?", customerId, userId).Updates(&cnvUpd)
	if qry.RowsAffected <= 0 {
		log.Println("\tupdate customer query error: data not found")
		return customer.Core{}, errors.New("not found")
	}

	if err := qry.Error; err != nil {
		log.Println("\tupdate customer query error: ", err.Error())
		return customer.Core{}, errors.New("not found")
	}

	return ToCore(cnvUpd), nil

}
