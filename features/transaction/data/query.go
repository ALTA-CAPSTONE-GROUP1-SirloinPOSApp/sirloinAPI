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
	"github.com/midtrans/midtrans-go/snap"
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

func (tq *transactionQuery) TotalPrice(uCart transaction.Cart) float64 {
	totalPrice := 0.0
	for _, val := range uCart.Items {
		totalPrice += float64(val.Quantity) * val.Price
	}
	return totalPrice
}

func (tq *transactionQuery) Discount(uCart transaction.Cart, totalPrice float64) (float64, float64) {
	disc := 0.0
	totalBill := totalPrice
	if uCart.CustomerId != 0 {
		disc = 0.10
		totalBill = totalPrice - (totalPrice * disc)
	}
	return totalBill, disc
}

func (tq *transactionQuery) CreateTransaction(userId uint, uCart transaction.Cart, productStatus string, totalPrice, disc, totalBill float64) Transaction {
	return Transaction{
		UserId:        userId,
		CustomerId:    uCart.CustomerId,
		TotalPrice:    totalPrice,
		Discount:      disc,
		TotalBill:     totalBill,
		CreatedAt:     time.Now(),
		ProductStatus: productStatus,
		OrderStatus:   "pending",
	}
}

func (tq *transactionQuery) CreateNumberInvoice(transInput Transaction) string {
	cnvID := strconv.Itoa(int(transInput.ID))
	invNo := fmt.Sprintf("INV/" + transInput.CreatedAt.Format("20060102") + "/" + cnvID)
	return invNo
}

func (tq *transactionQuery) CreateTransProds(transInput Transaction, uCart transaction.Cart) []TransactionProduct {
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
	return transProds
}

func (tq *transactionQuery) AddSell(userId uint, uCart transaction.Cart) (transaction.Core, error) {
	tx := tq.db.Begin()

	//menghitung total price
	totalPrice := tq.TotalPrice(uCart)

	//diskon customer terdaftar
	totalBill, disc := tq.Discount(uCart, totalPrice)

	//mmebuat transaksi
	transInput := tq.CreateTransaction(userId, uCart, "sell", totalPrice, disc, totalBill)

	//input transaksi ke tabel
	if err := tx.Create(&transInput).Error; err != nil {
		tx.Rollback()
		log.Println("error add order query: ", err.Error())
		return transaction.Core{}, err
	}

	//membuat nomor invoice
	invNo := tq.CreateNumberInvoice(transInput)

	//save no invoice ke tabel
	transInput.InvoiceNumber = invNo
	tx.Save(&transInput)

	//membuat transactionproduct
	transProds := tq.CreateTransProds(transInput, uCart)
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

func (tq *transactionQuery) AddBuy(userId uint, uCart transaction.Cart) (transaction.Core, error) {
	tx := tq.db.Begin()

	//menghitung total price
	totalPrice := tq.TotalPrice(uCart)

	//diskon customer terdaftar
	totalBill, disc := tq.Discount(uCart, totalPrice)

	//mmebuat transaksi
	transInput := tq.CreateTransaction(userId, uCart, "buy", totalPrice, disc, totalBill)

	//input transaksi ke tabel
	if err := tx.Create(&transInput).Error; err != nil {
		tx.Rollback()
		log.Println("error add order query: ", err.Error())
		return transaction.Core{}, err
	}

	//membuat nomor invoice
	invNo := tq.CreateNumberInvoice(transInput)

	//save no invoice ke tabel
	transInput.InvoiceNumber = invNo
	tx.Save(&transInput)

	//membuat transactionproduct
	transProds := tq.CreateTransProds(transInput, uCart)
	if err := tx.Create(&transProds).Error; err != nil {
		tx.Rollback()
		log.Println("error create transactionproduct: ", err.Error())
		return transaction.Core{}, err
	}

	// membuat pembayaran midtrans
	s := config.MidtransSnapClient()
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  transInput.InvoiceNumber,
			GrossAmt: int64(totalPrice),
		},
	}
	snapResp, err := s.CreateTransaction(req)
	if err != nil {
		tx.Rollback()
		log.Println("error making midtrans transaction: ", err.Error())
		return transaction.Core{}, err
	}

	//save payment url ke table
	transInput.PaymentUrl = snapResp.RedirectURL
	tx.Save(&transInput)

	// commit tx transaksi
	tx.Commit()

	return transaction.Core{}, nil
}
