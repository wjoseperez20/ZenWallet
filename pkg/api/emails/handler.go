package emails

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/zenwallet/pkg/gmail"
	"gopkg.in/gomail.v2"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
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

	if err := sendEmail(username); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// sendEmail sends an email with the account statement
func sendEmail(username string) error {

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
	accountStatementData := map[string]interface{}{
		"Username":      username,
		"TotalBalance":  200.00,
		"AverageDebit":  -206.00,
		"DebitCount":    4,
		"AverageCredit": 366.00,
		"CreditCount":   12,
		"TransactionCountByMonth": map[string]int{
			"January 2023":  5,
			"February 2023": 8,
			"March 2023":    12,
		},
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
