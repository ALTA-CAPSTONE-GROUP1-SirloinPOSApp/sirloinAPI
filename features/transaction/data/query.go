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
		UserId:            userId,
		CustomerId:        uCart.CustomerId,
		TotalPrice:        totalPrice,
		Discount:          disc,
		TotalBill:         totalBill,
		CreatedAt:         time.Now(),
		ProductStatus:     productStatus,
		TransactionStatus: "pending",
	}
}

// INV/tanggaltransaksi/(buy/sell)/idtransaction
func (tq *transactionQuery) CreateNumberInvoice(transInput Transaction, productStatus string) string {
	cnvID := strconv.Itoa(int(transInput.ID))
	invNo := fmt.Sprintf("INV/" + transInput.CreatedAt.Format("20060102") + "/" + productStatus + "/" + cnvID)
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
	productStatus := "sell"
	//menghitung total price
	totalPrice := tq.TotalPrice(uCart)

	//diskon customer terdaftar
	totalBill, disc := tq.Discount(uCart, totalPrice)

	//mmebuat transaksi
	transInput := tq.CreateTransaction(userId, uCart, productStatus, totalPrice, disc, totalBill)

	//input transaksi ke tabel
	if err := tx.Create(&transInput).Error; err != nil {
		tx.Rollback()
		log.Println("error add order query: ", err.Error())
		return transaction.Core{}, err
	}

	//membuat nomor invoice
	invNo := tq.CreateNumberInvoice(transInput, productStatus)

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
	productStatus := "buy"
	//menghitung total price
	totalPrice := tq.TotalPrice(uCart)

	//diskon customer terdaftar
	totalBill, disc := tq.Discount(uCart, totalPrice)

	//mmebuat transaksi
	transInput := tq.CreateTransaction(userId, uCart, productStatus, totalPrice, disc, totalBill)

	//input transaksi ke tabel
	if err := tx.Create(&transInput).Error; err != nil {
		tx.Rollback()
		log.Println("error add order query: ", err.Error())
		return transaction.Core{}, err
	}

	//membuat nomor invoice
	invNo := tq.CreateNumberInvoice(transInput, productStatus)

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
			GrossAmt: int64(totalBill),
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

	return DataToCoreT(transInput), nil
}

func (tq *transactionQuery) GetTransactionHistory(userId uint, status, from, to string) ([]transaction.Core, error) {
	trans := []transaction.Core{}

	var err error
	if from == "" && to == "" {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ?", userId, status).Scan(&trans).Error
	} else if to == "" {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ? AND t.created_at >= ?", userId, status, from).Scan(&trans).Error
	} else if from == "" {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ? AND t.created_at <= ?", userId, status, to).Scan(&trans).Error
	} else {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ? AND t.created_at >= ? AND t.created_at <= ?", userId, status, from, to).Scan(&trans).Error
	}
	if err != nil {
		log.Println("error query get transactions history: ", err)
		return []transaction.Core{}, err
	}

	return trans, nil
}

func (tq *transactionQuery) GetTransactionDetails(transactionId uint) (transaction.TransactionRes, error) {
	trans := transaction.Core{}

	err := tq.db.Raw("SELECT t.id , t.customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url  FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.id = ?", transactionId).Scan(&trans).Error
	if err != nil {
		log.Println("error select transaction: ", err.Error())
		return transaction.TransactionRes{}, err
	}

	transR := transaction.CoreToResp(trans)
	tp := []transaction.TransactionProductRes{}

	err = tq.db.Raw("SELECT tp.product_id , p.product_name , tp.price , quantity , tp.total_price , p.product_image FROM transaction_products tp JOIN products p ON tp.product_id = p.id WHERE transaction_id = ?", transactionId).Scan(&tp).Error
	if err != nil {
		log.Println("error select transaction_product: ", err.Error())
		return transaction.TransactionRes{}, err
	}

	transR.TransactionProductRes = append(transR.TransactionProductRes, tp...)

	return transR, nil
}
func (tq *transactionQuery) GetAdminTransactionHistory(status, from, to string) ([]transaction.AdmTransactionRes, error) {
	trans := []transaction.AdmTransactionRes{}

	var err error
	if from == "" && to == "" {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ?", status).Scan(&trans).Error
	} else if to == "" {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ? AND t.created_at >= ?", status, from).Scan(&trans).Error
	} else if from == "" {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ? AND t.created_at <= ?", status, to).Scan(&trans).Error
	} else {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ? AND t.created_at >= ? AND t.created_at <= ?", status, from, to).Scan(&trans).Error
	}
	if err != nil {
		log.Println("error query get transactions history: ", err)
		return []transaction.AdmTransactionRes{}, err
	}

	return trans, nil
}
func (tq *transactionQuery) GetAdminTransactionDetails(transactionId uint) (transaction.AdmTransactionResDet, error) {
	trans := transaction.AdmTransactionRes{}

	err := tq.db.Raw("SELECT t.id , t.customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , invoice_number , invoice_url , payment_url  FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.id = ?", transactionId).Scan(&trans).Error
	if err != nil {
		log.Println("error select transaction: ", err.Error())
		return transaction.AdmTransactionResDet{}, err
	}

	transR := transaction.ToAdmResp(trans)
	tp := []transaction.TransactionProductRes{}

	err = tq.db.Raw("SELECT tp.product_id , p.product_name , tp.price , quantity , tp.total_price , p.product_image FROM transaction_products tp JOIN products p ON tp.product_id = p.id WHERE transaction_id = ?", transactionId).Scan(&tp).Error
	if err != nil {
		log.Println("error select transaction_product: ", err.Error())
		return transaction.AdmTransactionResDet{}, err
	}

	transR.TransactionProductRes = append(transR.TransactionProductRes, tp...)

	return transR, nil
}
