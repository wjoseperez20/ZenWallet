package emails

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/zenwallet/pkg/gmail"
	"github.com/wjoseperez20/zenwallet/pkg/models"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// @BasePath /api/v1

// SendAccountStatementEmail godoc
// @Summary Send account statement by Email
// @Description
// @Tags Emails
// @Security JwtAuth
// @Accept  json
// @Success 200 {string} string "Email sent successfully"
// @Router /emails/{emails} [post]
func SendAccountStatementEmail(c *gin.Context) {

	username := "John Doe"
	message := "This is a sample email message."
	transactions := []models.Transaction{
		{Date: time.Time{}, Account: 10001, Amount: -50.00},
		{Date: time.Time{}, Account: 10001, Amount: -25.00},
		// Add more transactions as needed
	}

	if err := sendEmail(username, message, transactions); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// sendEmail sends an email with the account statement
func sendEmail(username, message string, transaction []models.Transaction) error {

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Parse the HTML template
	accountStatementPath := filepath.Join(basePath, "assets/account_statement_template.html")
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

	// Send the email
	if err := gmail.Mailer.DialAndSend(m); err != nil {
		log.Printf("error sending email: %v", err)
	}

	return nil
}
