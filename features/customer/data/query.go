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

func (cq *customerQry) GetUserCustomers(userId uint) ([]customer.Core, error) {
	res := []Customer{}
	if err := cq.db.Where("user_id = ?", userId).Find(&res).Error; err != nil {
		log.Println("\terror query get all customers: ", err.Error())
		return []customer.Core{}, errors.New("not found")
	}
	return ToCoreArr(res), nil
}

func (cq *customerQry) GetCustomerById(userId, customerId uint) (customer.Core, error) {
	res := Customer{}
	if err := cq.db.Where("id = ? AND user_id = ?", customerId, userId).First(&res).Error; err != nil {
		log.Println("\terror query get user customer:", err.Error())
		return customer.Core{}, errors.New("not found")
	}
	return ToCore(res), nil
}
