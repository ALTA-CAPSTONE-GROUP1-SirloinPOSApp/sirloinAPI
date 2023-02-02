package data

import (
	"fmt"
	"log"
	"sirloinapi/config"
	"sirloinapi/features/transaction"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"gorm.io/gorm"
)

type transactionQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) transaction.TransactionData {
	return &transactionQuery{
		db: db,
	}
}

func (tq *transactionQuery) AddSell(userId uint, uCart transaction.Cart) (transaction.Core, error) {
	tx := tq.db.Begin()

	//menghitung total price
	totalPrice := 0.0
	for _, val := range uCart.Items {
		totalPrice += float64(val.Quantity) * val.Price
	}

	//diskon customer terdaftar
	disc := 0.0
	totalBill := totalPrice
	if uCart.CustomerId != 0 {
		disc = 0.10
		totalBill = totalPrice - (totalPrice * disc)
	}

	//mmebuat transaksi
	transInput := Transaction{
		UserId:        userId,
		CustomerId:    uCart.CustomerId,
		TotalPrice:    totalPrice,
		Discount:      disc,
		TotalBill:     totalBill,
		CreatedAt:     time.Now(),
		ProductStatus: "sell",
		OrderStatus:   "pending",
	}

	//input transaksi ke tabel
	if err := tx.Create(&transInput).Error; err != nil {
		tx.Rollback()
		log.Println("error add order query: ", err.Error())
		return transaction.Core{}, err
	}

	//membuat nomor invoice
	cnvID := strconv.Itoa(int(transInput.ID))
	invNo := fmt.Sprintf("INV/" + transInput.CreatedAt.Format("20060102") + "/" + cnvID)

	//save no invoice ke tabel
	transInput.InvoiceNumber = invNo
	tx.Save(&transInput)

	//membuat transactionproduct
	transProds := []TransactionProduct{}
	for _, item := range uCart.Items {
		transProd := TransactionProduct{
			TransactionId: transInput.ID,
			ProductId:     item.ProductId,
			Quantity:      item.Quantity,
			Price:         item.Price,
			TotalPrice:    item.Price * float64(item.Quantity),
		}
		transProds = append(transProds, transProd)
	}
	if err := tx.Create(&transProds).Error; err != nil {
		tx.Rollback()
		log.Println("error create transactionproduct: ", err.Error())
		return transaction.Core{}, err
	}

	//membuat qrcode midtrans jika cashless
	if uCart.PaymentMethod == "cashless" {
		c := config.MidtransCoreAPIClient()

		resp, err := c.ChargeTransaction(&coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeQris,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  transInput.InvoiceNumber,
				GrossAmt: int64(transInput.TotalBill),
			}, Qris: &coreapi.QrisDetails{Acquirer: "airpay shopee"},
		})
		if err != nil {
			tx.Rollback()
			log.Println("error create midtrans transaction: ", err.Error())
			return transaction.Core{}, err
		}

		//save qr code url ke table
		transInput.PaymentUrl = resp.Actions[0].URL
		tx.Save(&transInput)
	}

	//commite tx
	tx.Commit()

	return DataToCoreT(transInput), nil
}
