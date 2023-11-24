package api

import (
	"github.com/wjoseperez20/zenwallet/docs"
	"github.com/wjoseperez20/zenwallet/pkg/api/accounts"
	"github.com/wjoseperez20/zenwallet/pkg/api/healtcheck"
	"github.com/wjoseperez20/zenwallet/pkg/auth"
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
		v1.GET("/", healtcheck.Healthcheck)
		v1.POST("/login", middleware.APIKeyAuth(), auth.LoginHandler)
		v1.POST("/register", middleware.APIKeyAuth(), auth.RegisterHandler)

		account := v1.Group("/accounts")
		{
			account.GET("/", middleware.JWTAuth(), accounts.FindAccounts)
			account.POST("/", middleware.JWTAuth(), accounts.CreateAccount)
			account.GET("/:id", middleware.JWTAuth(), accounts.FindAccount)
			account.PUT(":id", middleware.JWTAuth(), accounts.UpdateAccount)
			account.DELETE(":id", middleware.JWTAuth(), accounts.DeleteAccount)
		}
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return r
}
