package emails

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @BasePath /api/v1

// SendEmail godoc
// @Summary Send account statement by Email
// @Description
// @Tags Emails
// @Security JwtAuth
// @Accept  json
// @Success 200 {string} string "Email sent successfully"
// @Router /emails/{email} [post]
func SendEmail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Email sent successfully"})
}

// SendEmails godoc
// @Summary Send account to All accounts
// @Description
// @Tags Emails
// @Security JWTAuth
// @Accept  json
// @Success 200 {string} string "Emails sent successfully"
// @Router /emails [post]
func SendEmails(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Emails sent successfully"})
}
