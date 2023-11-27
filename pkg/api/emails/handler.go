package emails

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/zenwallet/pkg/database"
	"github.com/wjoseperez20/zenwallet/pkg/gmail"
	"github.com/wjoseperez20/zenwallet/pkg/models"
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
	var input models.Email
	var account models.Account
	var transactions []models.Transaction

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get the account by email
	if err := database.DB.Where("email = ?", input.Email).First(&account).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "email not found"})
		return
	}

	// Get All transactions by account
	if err := database.DB.Where("account_id = ?", account.Account).Find(&transactions).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "transactions not found"})
		return
	}

	// Send the email
	if err := sendEmail(account, transactions); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// sendEmail sends an email with the account statement
func sendEmail(account models.Account, transactions []models.Transaction) error {

	basePath, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// Calculate Statement Account
	transactionsByMonth := make(map[string]int)
	var totalDebit float32
	var countDebit int
	var totalCredit float32
	var countCredit int
	for _, transaction := range transactions {
		month := transaction.CreatedAt.Format("January 2006")
		transactionsByMonth[month]++

		if transaction.Amount < 0 {
			totalDebit += transaction.Amount
			countDebit++
		} else {
			totalCredit += transaction.Amount
			countCredit++
		}
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
		"Username":                account.Client,
		"TotalBalance":            account.Balance,
		"AverageDebit":            totalDebit,
		"DebitCount":              countDebit,
		"AverageCredit":           totalCredit,
		"CreditCount":             countCredit,
		"TransactionCountByMonth": transactionsByMonth,
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "zenwallet.app@gmail.com")
	m.SetHeader("To", account.Email)
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
