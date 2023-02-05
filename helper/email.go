package helper

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"sirloinapi/config"
)

func SendEmail(to, subject, body, filename string) error {
	from := "sirloinpos@gmail.com"
	appPassword := config.GMAILAPPPASSWORD
	smtpServer := "smtp.gmail.com:587"

	// Read the PDF file
	pdf := readFile(filename)
	encoded := base64.StdEncoding.EncodeToString(pdf)

	// Prepare headers
	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = `multipart/mixed; boundary="NextPart"`

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n--NextPart\r\n"
	message += "Content-Type: text/plain\r\n"
	message += "Content-Transfer-Encoding: base64\r\n"
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body)) + "\r\n"
	message += "--NextPart\r\n"
	message += "Content-Type: application/pdf\r\n"
	message += "Content-Disposition: attachment; filename=invoice.pdf\r\n"
	message += "Content-Transfer-Encoding: base64\r\n"
	message += "\r\n" + encoded + "\r\n"
	message += "--NextPart--"

	auth := smtp.PlainAuth("", from, appPassword, "smtp.gmail.com")
	err := smtp.SendMail(smtpServer, auth, from, []string{to}, []byte(message))
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

func readFile(filename string) []byte {
	pdf, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read file: %v\n", err)
		return nil
	}
	return pdf
}
