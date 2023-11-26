package api

import (
	"github.com/wjoseperez20/zenwallet/docs"
	"github.com/wjoseperez20/zenwallet/pkg/api/accounts"
	"github.com/wjoseperez20/zenwallet/pkg/api/emails"
	"github.com/wjoseperez20/zenwallet/pkg/api/files"
	"github.com/wjoseperez20/zenwallet/pkg/api/healtcheck"
	"github.com/wjoseperez20/zenwallet/pkg/api/transactions"
	"github.com/wjoseperez20/zenwallet/pkg/api/users"
	"github.com/wjoseperez20/zenwallet/pkg/middleware"
	"time"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func InitRouter() *gin.Engine {
	r := gin.Default()

	r.Use(gin.Logger())
	if gin.Mode() == gin.ReleaseMode {
		r.Use(middleware.Security())
		r.Use(middleware.Xss())
	}
	r.Use(middleware.Cors())
	r.Use(middleware.RateLimiter(rate.Every(1*time.Minute), 60)) // 60 requests per minute

	docs.SwaggerInfo.BasePath = "/api/v1"
	v1 := r.Group("/api/v1")
	{
		v1.GET("/_", healtcheck.Healthcheck)
		v1.POST("/login", middleware.APIKeyAuth(), users.LoginUser)
		v1.POST("/register", middleware.APIKeyAuth(), users.RegisterUser)

		account := v1.Group("/accounts")
		{
			account.GET("/", middleware.JWTAuth(), accounts.FindAccounts)
			account.GET("/:account", middleware.JWTAuth(), accounts.FindAccount)
			account.POST("/", middleware.JWTAuth(), accounts.CreateAccount)
			account.PUT("/:account", middleware.JWTAuth(), accounts.UpdateAccount)
			account.DELETE("/:account", middleware.JWTAuth(), accounts.DeleteAccount)
		}

		transaction := v1.Group("/transactions")
		{
			transaction.GET("/", middleware.JWTAuth(), transactions.FindTransactions)
			transaction.GET("/:id", middleware.JWTAuth(), transactions.FindTransaction)
			transaction.POST("/", middleware.JWTAuth(), transactions.CreateTransaction)
			transaction.PUT("/:id", middleware.JWTAuth(), transactions.UpdateTransaction)
			transaction.DELETE("/:id", middleware.JWTAuth(), transactions.DeleteTransaction)
		}

		file := v1.Group("/files")
		{
			file.GET("/", middleware.JWTAuth(), files.FindFiles)
			file.GET("/:id", middleware.JWTAuth(), files.FindFile)
			file.POST("/", middleware.JWTAuth(), files.UploadFile)
		}

		email := v1.Group("/emails")
		{
			email.POST("/", middleware.JWTAuth(), emails.SendAccountStatementEmail)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
