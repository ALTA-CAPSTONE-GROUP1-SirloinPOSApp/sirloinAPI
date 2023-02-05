package helper

import (
	"fmt"
	"sirloinapi/features/transaction"
	"strconv"

	"github.com/jung-kurt/gofpdf"
)

func GeneratePDF(reports []transaction.Core, filename string) error {

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
	if allSell && !allBuy {
		judul = "Penjualan"
	} else if allBuy && !allSell {
		judul = "Pembelian"
	} else {
		judul = "Penjualan & Pembelian"
	}
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Laporan "+judul)
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Detail %s:", judul))
	pdf.Ln(10)
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(8, 8, "No.")
	pdf.Cell(25, 8, "Tanggal")
	pdf.Cell(45, 8, "No. Transaksi")
	// pdf.Cell(45, 8, "Pelanggan")
	pdf.Cell(30, 8, "Total Harga")
	pdf.Cell(25, 8, "Diskon")
	pdf.Cell(25, 8, "Total Tagihan")
	pdf.Cell(30, 8, "Status Transaksi")
	pdf.Ln(8)

	// Add data to the table
	for i, v := range reports {
		pdf.Cell(8, 8, fmt.Sprintf("%d", i+1))
		pdf.Cell(25, 8, fmt.Sprintf("%d/%d/%d", v.CreatedAt.Day(), v.CreatedAt.Month(), v.CreatedAt.Year()))
		pdf.Cell(45, 8, v.InvoiceNumber)
		// pdf.Cell(45, 8, v.CustomerName)
		totalPrice := v.TotalPrice
		totalPriceString := strconv.FormatFloat(totalPrice, 'f', 2, 64)
		pdf.Cell(30, 8, totalPriceString)
		discount := v.Discount
		discountString := strconv.FormatFloat(discount, 'f', 2, 64)
		pdf.Cell(25, 8, discountString)
		totalBill := v.TotalBill
		totalBillString := strconv.FormatFloat(totalBill, 'f', 2, 64)
		pdf.Cell(25, 8, totalBillString)
		pdf.Cell(30, 8, v.TransactionStatus)
		pdf.Ln(8)
	}
	if err := pdf.OutputFileAndClose(filename + "laporan.pdf"); err != nil {
		println(err)
		return err
	}

	return nil
}
