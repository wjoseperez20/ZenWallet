package emails

import (
	"github.com/gin-gonic/gin"
	"github.com/wjoseperez20/zenwallet/pkg/models"
	"github.com/wjoseperez20/zenwallet/pkg/smtp"
	"net/http"
	"time"
)

type EmailTemplateData struct {
	Username string
	Message  string
}

// @BasePath /api/v1

// AccountStatementEmail godoc
// @Summary Send account statement by Email
// @Description
// @Tags Emails
// @Security JwtAuth
// @Accept  json
// @Success 200 {string} string "Email sent successfully"
// @Router /emails/{emails} [post]
func AccountStatementEmail(c *gin.Context) {

	username := "John Doe"
	message := "This is a sample email message."
	transactions := []models.Transaction{
		{Date: time.Time{}, Account: 10001, Amount: -50.00},
		{Date: time.Time{}, Account: 10001, Amount: -25.00},
		// Add more transactions as needed
	}

	if err := smtp.SendEmail(username, message, transactions); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}
