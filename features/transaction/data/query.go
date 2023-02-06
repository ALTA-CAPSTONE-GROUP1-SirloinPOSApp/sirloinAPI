package data

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sirloinapi/config"
	product "sirloinapi/features/product/data"
	"sirloinapi/features/transaction"
	"sirloinapi/helper"
	"strconv"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/jung-kurt/gofpdf"
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

func (tq *transactionQuery) CheckAllowedProd(uid uint, uCart transaction.Cart) bool {
	check := true
	for _, val := range uCart.Items {
		p := product.Product{}
		tq.db.First(&p, val.ProductId)
		if p.UserId != uid {
			check = false
		}
	}
	return check
}

func (tq *transactionQuery) TotalPrice(uCart transaction.Cart) float64 {
	totalPrice := 0.0
	for _, val := range uCart.Items {
		p := product.Product{}
		tq.db.First(&p, val.ProductId)
		// totalPrice += float64(val.Quantity) * val.Price
		totalPrice += float64(val.Quantity) * p.Price
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

// INV-tanggaltransaksi-(buy/sell)-idtransaction
func (tq *transactionQuery) CreateNumberInvoice(transInput Transaction, productStatus string) string {
	cnvID := strconv.Itoa(int(transInput.ID))
	invNo := fmt.Sprintf("INV-" + transInput.CreatedAt.Format("20060102") + "-" + strings.ToUpper(productStatus) + "-" + cnvID)
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

// menambahkan 1 hari di query param to pada history laporan
func AddOneDay(to string) string {
	// Ganti karakter penghubung "/" menjadi "-"
	newDateString := strings.Replace(to, "/", "-", -1)

	// Parsing string menjadi time.Time
	t, _ := time.Parse("2006-01-02", newDateString)
	fmt.Println("Current time:", t)

	// Tambahkan 1 hari
	t = t.Add(time.Hour * 24)
	return t.String()
}

func (tq *transactionQuery) AddSell(userId uint, uCart transaction.Cart) (transaction.Core, error) {
	//check if the product really the seller product
	check := tq.CheckAllowedProd(userId, uCart)
	if !check {
		log.Println("unauthorized: request body contain product that's not belong to the user")
		return transaction.Core{}, errors.New("unauthorized: request body contain product that's not belong to the user")
	}

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
	check := tq.CheckAllowedProd(uint(1), uCart)
	if !check {
		log.Println("unauthorized: request body contain product that's not belong to super admin")
		return transaction.Core{}, errors.New("unauthorized: request body contain product that's not belong to super admin")
	}

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

	// add one day for query param 'to'
	to = AddOneDay(to)
	var err error
	if from == "" && to == "" {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ?", userId, status).Scan(&trans).Error
	} else if to == "" {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ? AND t.created_at >= ?", userId, status, from).Scan(&trans).Error
	} else if from == "" {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ? AND t.created_at <= ?", userId, status, to).Scan(&trans).Error
	} else {
		err = tq.db.Raw("SELECT t.id , c.id customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.user_id = ? AND product_status = ? AND t.created_at >= ? AND t.created_at <= ?", userId, status, from, to).Scan(&trans).Error
	}
	if err != nil {
		log.Println("error query get transactions history: ", err)
		return []transaction.Core{}, err
	}

	return trans, nil
}

func (tq *transactionQuery) GetTransactionDetails(transactionId uint) (transaction.TransactionRes, error) {
	trans := transaction.Core{}

	err := tq.db.Raw("SELECT t.id , t.customer_id , c.name customer_name , total_price , discount , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url  FROM transactions t JOIN customers c ON t.customer_id = c.id WHERE t.id = ?", transactionId).Scan(&trans).Error
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
	// Tambahkan 1 hari
	to = AddOneDay(to)
	var err error
	if from == "" && to == "" {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ?", status).Scan(&trans).Error
	} else if to == "" {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ? AND t.created_at >= ?", status, from).Scan(&trans).Error
	} else if from == "" {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ? AND t.created_at <= ?", status, to).Scan(&trans).Error
	} else {
		err = tq.db.Raw("SELECT t.id , user_id tenant_id , business_name tenant_name , total_bill , t.created_at , transaction_status , product_status , invoice_number , invoice_url , payment_url FROM transactions t JOIN users u ON u.id = t.user_id WHERE product_status = ? AND t.created_at >= ? AND t.created_at <= ?", status, from, to).Scan(&trans).Error
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
func (tq *transactionQuery) NotificationTransactionStatus(invNo, transStatus string) error {
	trans := Transaction{}

	err := tq.db.First(&trans, "invoice_number = ?", invNo).Error
	if err != nil {
		log.Println("error select transaction: ", err.Error())
		return err
	}

	// 5. Do set transaction status based on response from check transaction status
	if transStatus == "capture" {
		if transStatus == "challenge" {
			// TODO set transaction status on your database to 'challenge'
			// e.g: 'Payment status challenged. Please take action on your Merchant Administration Portal
			trans.TransactionStatus = "challenge"
		} else if transStatus == "accept" {
			// TODO set transaction status on your database to 'success'
			trans.TransactionStatus = "success"
		}
	} else if transStatus == "settlement" {
		// TODO set transaction status on your databaase to 'success'
		trans.TransactionStatus = "success"
	} else if transStatus == "cancel" || transStatus == "expire" {
		// TODO set transaction status on your databaase to 'failure'
		trans.TransactionStatus = "failure"
	} else if transStatus == "pending" {
		// TODO set transaction status on your databaase to 'pending' / waiting payment
		trans.TransactionStatus = "waiting payment"
	} else {
		trans.TransactionStatus = transStatus
	}

	aff := tq.db.Save(&trans)
	if aff.RowsAffected <= 0 {
		log.Println("error update transaction status, no rows affected")
		return errors.New("error update transaction status")
	}

	//update stock product
	if trans.TransactionStatus == "success" {
		transProds := []TransactionProduct{}
		tq.db.Find(&transProds, "transaction_id", trans.ID)
		for _, item := range transProds {
			prod := product.Product{}
			tq.db.First(&prod, item.ProductId)
			prod.Stock -= item.Quantity
			prod.ItemsSold += item.Quantity
			tq.db.Save(&prod)
		}

		if trans.ProductStatus == "sell" {
			if trans.CustomerId != uint(0) {
				//bikin invoice penjualan, upload ke s3 dan kirim email
				invURL, err := tq.Invoice(trans.Discount, trans.ID, true, trans.ProductStatus)
				if err != nil {
					return err
				}
				trans.InvoiceUrl = invURL
				tq.db.Save(&trans)
			} else {
				//bikin invoice penjualan dan upload ke s3
				invURL, err := tq.Invoice(trans.Discount, trans.ID, false, trans.ProductStatus)
				if err != nil {
					return err
				}
				trans.InvoiceUrl = invURL
				tq.db.Save(&trans)
			}
		} else if trans.ProductStatus == "buy" {
			//bikin invoice pembelian
			invURL, err := tq.Invoice(trans.Discount, trans.ID, false, trans.ProductStatus)
			if err != nil {
				return err
			}
			trans.InvoiceUrl = invURL
			tq.db.Save(&trans)
		}
	}

	return nil
}
func (tq *transactionQuery) UpdateStatus(transId uint, status string) error {
	input := Transaction{}
	err := tq.db.Where("id = ?", transId).First(&input).Error
	if err != nil {
		log.Println("error select transaction: ", err.Error())
		return err
	}
	input.TransactionStatus = status
	err = tq.db.Save(&input).Error
	if err != nil {
		log.Println("error save transaction status: ", err.Error())
		return err
	} else {
		//update stock product
		if input.TransactionStatus == "success" {
			transProds := []TransactionProduct{}
			tq.db.Find(&transProds, "transaction_id", input.ID)
			for _, item := range transProds {
				prod := product.Product{}
				tq.db.First(&prod, item.ProductId)
				prod.Stock -= item.Quantity
				prod.ItemsSold += item.Quantity
				tq.db.Save(&prod)
			}

			if input.ProductStatus == "sell" {
				if input.CustomerId != uint(0) {
					//bikin invoice penjualan, upload ke s3 dan kirim email
					invURL, err := tq.Invoice(input.Discount, input.ID, true, input.ProductStatus)
					if err != nil {
						return err
					}
					input.InvoiceUrl = invURL
					tq.db.Save(&input)
				} else {
					//bikin invoice penjualan dan upload ke s3
					invURL, err := tq.Invoice(input.Discount, input.ID, false, input.ProductStatus)
					if err != nil {
						return err
					}
					input.InvoiceUrl = invURL
					tq.db.Save(&input)
				}
			} else if input.ProductStatus == "buy" {
				//bikin invoice pembelian
				invURL, err := tq.Invoice(input.Discount, input.ID, false, input.ProductStatus)
				if err != nil {
					return err
				}
				input.InvoiceUrl = invURL
				tq.db.Save(&input)
			}
		}
	}
	return nil
}

func (tq *transactionQuery) Invoice(discount float64, transId uint, member bool, status string) (string, error) {
	tx := tq.db.Begin()

	transInv := transaction.TransactionInv{}

	err := tx.Raw("SELECT invoice_number , t.created_at transaction_date , u.business_name seller_name , u.email seller_email , u.phone_number seller_phone , u.address seller_address , c.name customer_name , c.email customer_email , c.phone_number customer_phone , c.address customer_address , t.total_price sub_total , discount , total_bill total_price FROM transactions t JOIN users u ON t.user_id = u.id JOIN customers c ON t.customer_id = c.id WHERE t.id = ?", transId).Scan(&transInv).Error
	if err != nil {
		tx.Rollback()
		log.Println("error select transaction invoice: ", err.Error())
		return "", err
	}

	tInv := transaction.InvToDetail(transInv)

	itms := []transaction.ItemsInv{}

	err = tq.db.Raw("SELECT p.product_name item_name , quantity , p.price , total_price FROM transaction_products tp JOIN products p ON tp.product_id = p.id WHERE tp.transaction_id = ?", transId).Scan(&itms).Error
	if err != nil {
		tx.Rollback()
		log.Println("error select transaction item invoice: ", err.Error())
		return "", err
	}

	tInv.Items = append(tInv.Items, itms...)

	tInv.DiscountAmount = tInv.SubTotal * discount

	//create new PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	// Add the header information
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "INVOICE", "0", 1, "C", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(100, 5, "", "0", 0, "C", false, 0, "")
	pdf.CellFormat(40, 5, "No. Invoice", "0", 0, "L", false, 0, "")
	pdf.CellFormat(50, 5, fmt.Sprint(tInv.InvoiceNumber), "0", 1, "C", false, 0, "")
	pdf.CellFormat(100, 5, "", "0", 0, "C", false, 0, "")
	pdf.CellFormat(40, 5, "Tanggal Transaksi:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(50, 5, fmt.Sprint(tInv.TransactionDate.Format("2006-01-02")), "0", 1, "C", false, 0, "")

	if status == "sell" {
		// Add the seller information
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(190, 5, "Diterbitkan oleh:", "0", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(190, 5, fmt.Sprint("Nama: \t"+tInv.SellerName), "0", 1, "L", false, 0, "")
		pdf.CellFormat(190, 5, fmt.Sprint("Telepon: \t"+tInv.SellerPhone), "0", 1, "L", false, 0, "")
		pdf.CellFormat(190, 5, fmt.Sprint("Alamat: \t"+tInv.SellerAddress), "0", 1, "L", false, 0, "")
	} else if status == "buy" {
		// Add the seller information
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(190, 5, "Kepada:", "0", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(190, 5, fmt.Sprint("Nama: \t"+tInv.SellerName), "0", 1, "L", false, 0, "")
		pdf.CellFormat(190, 5, fmt.Sprint("Telepon: \t"+tInv.SellerPhone), "0", 1, "L", false, 0, "")
		pdf.CellFormat(190, 5, fmt.Sprint("Alamat: \t"+tInv.SellerAddress), "0", 1, "L", false, 0, "")
	} else {
		log.Println("status empty string")
		return "", errors.New("bad request")
	}

	if member {
		// Add the customer information
		pdf.Ln(5)
		pdf.SetFont("Arial", "B", 10)
		pdf.CellFormat(190, 5, "Kepada:", "0", 1, "L", false, 0, "")
		pdf.SetFont("Arial", "", 10)
		pdf.CellFormat(190, 5, fmt.Sprint("Nama: \t"+tInv.CustomerName), "0", 1, "L", false, 0, "")
		pdf.CellFormat(190, 5, fmt.Sprint("Email: \t"+tInv.CustomerEmail), "0", 1, "L", false, 0, "")
		pdf.CellFormat(190, 5, fmt.Sprint("Telepon: \t"+tInv.CustomerPhone), "0", 1, "L", false, 0, "")
		pdf.CellFormat(190, 5, fmt.Sprint("Alamat: \t"+tInv.CustomerAddress), "0", 1, "L", false, 0, "")
	}

	// Add the item table
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(85, 5, "Nama Item", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 5, "Jumlah", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 5, "Harga", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 5, "Total Harga", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 10)
	for _, item := range tInv.Items {
		pdf.CellFormat(85, 5, item.ItemName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 5, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 5, fmt.Sprintf("Rp. %s", humanize.Commaf(item.Price)+",00"), "1", 0, "C", false, 0, "")
		pdf.CellFormat(40, 5, fmt.Sprintf("Rp. %s", humanize.Commaf(item.TotalPrice)+",00"), "1", 1, "C", false, 0, "")
	}

	// Add the subtotal, discount, and total price
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(100, 5, "", "0", 0, "C", false, 0, "")
	pdf.CellFormat(40, 5, "Subtotal:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(50, 5, fmt.Sprintf("Rp. %s", humanize.Commaf(tInv.SubTotal))+",00", "0", 1, "C", false, 0, "")
	pdf.CellFormat(100, 5, "", "0", 0, "C", false, 0, "")
	pdf.CellFormat(40, 5, "Diskon:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(50, 5, fmt.Sprintf("Rp. %s (%d%%)", humanize.Commaf(tInv.DiscountAmount)+",00", int(tInv.Discount*100)), "0", 1, "C", false, 0, "")
	pdf.CellFormat(100, 5, "", "0", 0, "C", false, 0, "")
	pdf.CellFormat(40, 5, "Total Pembayaran:", "0", 0, "L", false, 0, "")
	pdf.CellFormat(50, 5, fmt.Sprintf("Rp. %s", humanize.Commaf(tInv.TotalPrice)+",00"), "0", 1, "C", false, 0, "")

	// Save the PDF file
	pdf.OutputFileAndClose(fmt.Sprint(tInv.InvoiceNumber) + ".pdf")

	//send buying invoice email to tenant or to registered customer
	if status == "sell" && member {
		body := "Dear " + tInv.CustomerName + ",\nBerikut adalah invoice untuk transaksi dengan nomor: " + transInv.InvoiceNumber + "\n\nEmail ini dibuat secara otomatis, mohon untuk tidak membalas email ini. \n\nTerima Kasih"
		err := helper.SendEmail(tInv.CustomerEmail, "INVOICE", body, fmt.Sprint(tInv.InvoiceNumber)+".pdf")
		if err != nil {
			log.Println("error sending email to customer: ", err.Error())
			return "", err
		}
	} else if status == "buy" {
		body := "Dear " + tInv.SellerName + ",\n Berikut adalah invoice untuk transaksi dengan nomor: " + transInv.InvoiceNumber + "\n\nPesan ini dibuat secara otomatis, mohon untuk tidak membalas email ini. \n\nTerima Kasih"
		err := helper.SendEmail(tInv.SellerEmail, "INVOICE", body, fmt.Sprint(tInv.InvoiceNumber)+".pdf")
		if err != nil {
			log.Println("error sending email to tenant: ", err.Error())
			return "", err
		}
	}

	//read pdf
	file, err := os.Open(fmt.Sprint(tInv.InvoiceNumber) + ".pdf")
	if err != nil {
		log.Println("error open file: ", err)
		log.Fatal(err)
		return "", err
	}

	//upload pdf to s3
	res, err := helper.UploadPdfToS3(fmt.Sprint("/files/invoice/", tInv.InvoiceNumber, ".pdf"), file)
	if err != nil {
		log.Println("error upload pdf to s3: ", err.Error())
		return "", err
	}

	//delete pdf file on local
	err = os.Remove(fmt.Sprint(tInv.InvoiceNumber) + ".pdf")
	if err != nil {
		log.Println("error delete file: ", err.Error())
		return "", err
	}
	tx.Commit()
	return res, nil
}
