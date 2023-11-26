package main

import (
	"github.com/wjoseperez20/zenwallet/pkg/amazon"
	"github.com/wjoseperez20/zenwallet/pkg/api"
	"github.com/wjoseperez20/zenwallet/pkg/cache"
	"github.com/wjoseperez20/zenwallet/pkg/database"
	"github.com/wjoseperez20/zenwallet/pkg/gmail"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cache.InitRedis()
	database.ConnectDatabase()
	amazon.ConnectAWS()
	gmail.ConnectGmail()

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	r := api.InitRouter()

	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}
