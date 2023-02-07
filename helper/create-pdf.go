package helper

import (
	"fmt"
	"sirloinapi/features/transaction"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDFReport(data interface{}, filename string) error {
	reports := ToReports(data)
	allSell := true
	allBuy := true
	for _, v := range reports {
		if v.ProductStatus == "sell" {
			allBuy = false
		}
		if v.ProductStatus == "buy" {
			allSell = false
		}
	}
	judul := ""
	tenant := ""
	plusMinus := ""
	if allSell && !allBuy {
		judul = "Penjualan"
	} else if allBuy && !allSell {
		judul = "Belanja"
	} else {
		judul = "Penjualan & Belanja"
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	if reports[0].TenantName != "" {
		tenant = "Tenant " + reports[0].TenantName
	}
	pdf.Cell(40, 10, fmt.Sprintf("Laporan %s %s", judul, tenant))
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Detail %s:", judul))
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(8, 8, "No.")
	pdf.Cell(25, 8, "Tanggal")
	pdf.Cell(45, 8, "No. Transaksi")
	// pdf.Cell(45, 8, "Pelanggan")
	if judul != "Belanja" {
		pdf.Cell(30, 8, "Total Harga")
		pdf.Cell(25, 8, "Diskon")
		pdf.Cell(25, 8, "Total Tagihan")
	} else {
		pdf.Cell(25, 8, "Total Belanja")
	}

	pdf.Cell(30, 8, "Status Transaksi")
	pdf.Ln(8)

	// Add data to the table
	for i, v := range reports {
		if v.ProductStatus == "sell" && v.TransactionStatus == "success" {
			plusMinus = "+"
		} else if v.ProductStatus == "buy" && v.TransactionStatus == "success" {
			plusMinus = "-"
		} else {
			plusMinus = ""
		}
		pdf.Cell(8, 8, fmt.Sprintf("%d", i+1))
		pdf.Cell(25, 8, fmt.Sprintf("%d/%d/%d", v.CreatedAt.Day(), v.CreatedAt.Month(), v.CreatedAt.Year()))
		pdf.Cell(45, 8, v.InvoiceNumber)
		// pdf.Cell(45, 8, v.CustomerName)
		if judul != "Belanja" {
			totalPrice := v.TotalPrice
			totalPriceString := strconv.FormatFloat(totalPrice, 'f', 2, 64)
			pdf.Cell(30, 8, totalPriceString)
			discount := v.Discount
			discountString := strconv.FormatFloat(discount, 'f', 2, 64)
			pdf.Cell(25, 8, discountString)
		}
		totalBill := v.TotalBill
		totalBillString := strconv.FormatFloat(totalBill, 'f', 2, 64)
		pdf.Cell(25, 8, plusMinus+totalBillString)
		pdf.Cell(30, 8, v.TransactionStatus)
		pdf.Ln(8)
	}
	if err := pdf.OutputFileAndClose(filename + "laporan.pdf"); err != nil {
		return err
	}

	return nil
}

type Report struct {
	ProductStatus     string
	TenantName        string
	CreatedAt         time.Time
	InvoiceNumber     string
	TotalPrice        float64
	Discount          float64
	TotalBill         float64
	TransactionStatus string
}

func ToReports(data interface{}) []Report {
	res := []Report{}

	switch docs := data.(type) {
	case []transaction.Core:
		for _, v := range docs {
			tmp := CoreToReport(v)
			res = append(res, tmp)
		}
	case []transaction.AdmTransactionRes:
		for _, v := range docs {
			tmp := AdmToReport(v)
			res = append(res, tmp)
		}
	default:
		return nil
	}
	return res
}

func CoreToReport(data transaction.Core) Report {
	return Report{
		ProductStatus:     data.ProductStatus,
		TenantName:        data.TenantName,
		CreatedAt:         data.CreatedAt,
		InvoiceNumber:     data.InvoiceNumber,
		TotalPrice:        data.TotalPrice,
		Discount:          data.Discount,
		TotalBill:         data.TotalBill,
		TransactionStatus: data.TransactionStatus,
	}
}

func AdmToReport(data transaction.AdmTransactionRes) Report {
	return Report{
		ProductStatus:     "buy",
		TenantName:        data.TenantName,
		CreatedAt:         data.CreatedAt,
		InvoiceNumber:     data.InvoiceNumber,
		TotalPrice:        0.00,
		Discount:          0.00,
		TotalBill:         data.TotalBill,
		TransactionStatus: data.TransactionStatus,
	}
}
