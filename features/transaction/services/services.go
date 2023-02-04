package services

import (
	"errors"
	"log"
	"sirloinapi/features/transaction"
	"sirloinapi/helper"
	"strings"
)

type transSvc struct {
	qry transaction.TransactionData
}

func New(data transaction.TransactionData) transaction.TransactionService {
	return &transSvc{
		qry: data,
	}
}

func (ts *transSvc) AddSell(token interface{}, uCart transaction.Cart) (transaction.Core, error) {
	userId := helper.ExtractToken(token)
	if userId <= 0 {
		log.Println("error extract token add order")
		return transaction.Core{}, errors.New("user not found")
	}

	res, err := ts.qry.AddSell(uint(userId), uCart)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "bad request") {
			msg = "bad request"
		} else {
			msg = "server problem"
		}
		log.Println("error add order query in service : ", err.Error())
		return transaction.Core{}, errors.New(msg)
	}

	return res, nil
}
func (ts *transSvc) AddBuy(token interface{}, uCart transaction.Cart) (transaction.Core, error) {
	userId := helper.ExtractToken(token)
	if userId <= 0 {
		log.Println("error extract token add order")
		return transaction.Core{}, errors.New("user not found")
	}

	res, err := ts.qry.AddBuy(uint(userId), uCart)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "bad request") {
			msg = "bad request"
		} else {
			msg = "server problem"
		}
		log.Println("error add order query in service : ", err.Error())
		return transaction.Core{}, errors.New(msg)
	}

	return res, nil
}
func (ts *transSvc) GetTransactionHistory(token interface{}, status, from, to string) ([]transaction.Core, error) {
	userId := helper.ExtractToken(token)
	if userId <= 0 {
		log.Println("error extract token")
		return []transaction.Core{}, errors.New("user not found")
	}

	res, err := ts.qry.GetTransactionHistory(uint(userId), status, from, to)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		log.Println("error calling gettransactionhistory data in service: ", err.Error())
		return []transaction.Core{}, errors.New(msg)
	}

	return res, nil
}

func (ts *transSvc) GetTransactionDetails(transactionId uint) (transaction.TransactionRes, error) {
	res, err := ts.qry.GetTransactionDetails(transactionId)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		log.Println("error calling GetTransactionDetails data in service: ", err.Error())
		return transaction.TransactionRes{}, errors.New(msg)
	}

	return res, nil
}
func (ts *transSvc) GetAdminTransactionHistory(status, from, to string) ([]transaction.AdmTransactionRes, error) {
	res, err := ts.qry.GetAdminTransactionHistory(status, from, to)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		log.Println("error calling GetAdminTransactionHistory data in service: ", err.Error())
		return []transaction.AdmTransactionRes{}, errors.New(msg)
	}

	return res, nil
}
func (ts *transSvc) GetAdminTransactionDetails(transactionId uint) (transaction.AdmTransactionResDet, error) {
	res, err := ts.qry.GetAdminTransactionDetails(transactionId)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "data not found"
		} else {
			msg = "server problem"
		}
		log.Println("error calling GetTransactionDetails data in service: ", err.Error())
		return transaction.AdmTransactionResDet{}, errors.New(msg)
	}

	return res, nil
}
