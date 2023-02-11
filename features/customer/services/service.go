package services

import (
	"errors"
	"fmt"
	"log"
	"sirloinapi/features/customer"
	"sirloinapi/helper"
	"strings"

	"github.com/go-playground/validator/v10"
)

type customerUseCase struct {
	qry customer.CustomerData
	vld *validator.Validate
}

func New(cd customer.CustomerData) customer.CustomerService {
	return &customerUseCase{
		qry: cd,
		vld: validator.New(),
	}
}

func (cuc *customerUseCase) Add(userToken interface{}, newCustomer customer.Core) (customer.Core, error) {
	userId := helper.ExtractToken(userToken)
	if userId <= 1 {
		log.Println("extract token error, not allowed to access customer")
		return customer.Core{}, errors.New("extract token error, not allowed to access customer")
	}
	err := helper.Validasi(helper.ToValidate("customer", newCustomer))
	if err != nil {
		return customer.Core{}, err
	}
	res, err := cuc.qry.Add(uint(userId), newCustomer)
	if err != nil {
		errmsg := ""
		if strings.Contains(err.Error(), "not found") {
			errmsg = "data not found"
		} else if strings.Contains(err.Error(), "Duplicate") && strings.Contains(err.Error(), "customers.email") {
			errmsg = "user already exist"
		} else if strings.Contains(err.Error(), "Duplicate") && strings.Contains(err.Error(), "customers.phone_number") {
			errmsg = "phone number already exist"
		} else {
			errmsg = "server problem"
		}
		log.Println("error update query: ", err.Error())
		return customer.Core{}, errors.New(errmsg)
	}
	return res, nil
}

func (cuc *customerUseCase) Update(userToken interface{}, customerId uint, updateData customer.Core) (customer.Core, error) {
	userId := helper.ExtractToken(userToken)
	if userId <= 1 {
		log.Println("extract token error, not allowed to access customer")
		return customer.Core{}, errors.New("extract token error, not allowed to access customer")
	}
	if updateData.Name != "" {
		err := helper.Validasi(helper.ToValidate("as", updateData))
		if err != nil {
			return customer.Core{}, fmt.Errorf("update customer name: , %v", err)
		}
	}

	if updateData.Email != "" {
		err := helper.Validasi(helper.ToValidate("email", updateData))
		if err != nil {
			return customer.Core{}, fmt.Errorf("update customer email: , %v", err)
		}
	}

	if updateData.PhoneNumber != "" {
		err := helper.Validasi(helper.ToValidate("pn", updateData))
		if err != nil {
			return customer.Core{}, fmt.Errorf("update customer phone number: , %v", err)
		}
	}
	res, err := cuc.qry.Update(uint(userId), customerId, updateData)
	if err != nil {
		errmsg := ""
		if strings.Contains(err.Error(), "not found") {
			errmsg = "data not found"
		} else if strings.Contains(err.Error(), "Duplicate") && strings.Contains(err.Error(), "customers.email") {
			errmsg = "user already exist"
		} else if strings.Contains(err.Error(), "Duplicate") && strings.Contains(err.Error(), "customers.phone_number") {
			errmsg = "phone number already exist"
		} else {
			errmsg = "server problem"
		}
		log.Println("error update query: ", err.Error())
		return customer.Core{}, errors.New(errmsg)
	}
	return res, nil
}

func (cuc *customerUseCase) GetUserCustomers(token interface{}) ([]customer.Core, error) {
	userId := helper.ExtractToken(token)
	if userId <= 1 {
		log.Println("extract token error, not allowed to access customer")
		return []customer.Core{}, errors.New("extract token error, not allowed to access customer")
	}
	res, err := cuc.qry.GetUserCustomers(uint(userId))
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		return []customer.Core{}, errors.New(msg)
	}
	return res, nil
}

func (cuc *customerUseCase) GetCustomerById(token interface{}, customerId uint) (customer.Core, error) {
	userId := helper.ExtractToken(token)
	if userId <= 1 {
		log.Println("extract token error, not allowed to access customer")
		return customer.Core{}, errors.New("extract token error, not allowed to access customer")
	}
	res, err := cuc.qry.GetCustomerById(uint(userId), customerId)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		return customer.Core{}, errors.New(msg)
	}
	return res, nil
}
