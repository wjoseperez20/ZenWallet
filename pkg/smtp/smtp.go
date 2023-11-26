package smtp

import (
	"bytes"
	"github.com/wjoseperez20/zenwallet/pkg/models"
	"html/template"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/gomail.v2"
)

// SendEmail godoc
// @Summary Send account statement by Email
func SendEmail(username, message string, transaction []models.Transaction) error {

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the HTML template
	accountStatementPath := filepath.Join(basePath, "/pkg/smtp/account_statement_template.html")
	accountStatementTemplate, err := template.ParseFiles(accountStatementPath)
	if err != nil {
		log.Printf("failed to parse email template: %v", err)
		return nil
	}

	// Data for the emails template
	accountStatementData := models.EmailTemplateData{
		Username:     username,
		Message:      message,
		Transactions: transaction,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "zenwallet.app@gmail.com")
	m.SetHeader("To", "zenwallet.app@gmail.com")
	m.SetHeader("Subject", "Account Statement")

	// Render the template with the data
	var renderedBody bytes.Buffer
	err = accountStatementTemplate.Execute(&renderedBody, accountStatementData)
	if err != nil {
		log.Printf("failed to render email template: %v", err)
	}

	m.SetBody("text/html", renderedBody.String())

	// Set up the SMTP server details
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "zenwallet.app@gmail.com", os.Getenv("SMTP_SECRET"))
	dialer.TLSConfig = nil

	// Send the emails
	if err := dialer.DialAndSend(m); err != nil {
		log.Printf("error sending email: %v", err)
	}

	return nil
}
