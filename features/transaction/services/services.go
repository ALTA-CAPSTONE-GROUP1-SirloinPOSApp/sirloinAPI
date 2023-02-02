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
