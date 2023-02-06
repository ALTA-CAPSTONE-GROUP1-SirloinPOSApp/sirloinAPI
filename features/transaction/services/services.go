package services

import (
	"errors"
	"log"
	"sirloinapi/config"
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
func (ts *transSvc) GetTransactionHistory(token interface{}, status, from, to, sendEmail string) ([]transaction.Core, error) {
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
func (ts *transSvc) NotificationTransactionStatus(invNo string) error {
	c := config.MidtransCoreAPIClient()

	// 4. Check transaction to Midtrans with param invoice number
	transactionStatusResp, e := c.CheckTransaction(invNo)
	if e != nil {
		log.Println("error check transaction status: ", e.Error())
		return errors.New("error check transaction status")
	}

	err := ts.qry.NotificationTransactionStatus(invNo, transactionStatusResp.TransactionStatus)
	if err != nil {
		log.Println("error calling NotificationTransactionStatus data in service: ", err.Error())
		return errors.New("error calling NotificationTransactionStatus data in service")
	}

	return nil
}
func (ts *transSvc) UpdateStatus(transId uint, status string) error {
	err := ts.qry.UpdateStatus(transId, status)
	if err != nil {
		msg := ""
		if strings.Contains(err.Error(), "not found") {
			msg = "bad request"
		} else {
			msg = "server problem"
		}
		return errors.New(msg)
	}
	return nil
}
