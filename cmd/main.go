package main

import (
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/logging"
	"github.com/gin-gonic/gin"
)

func HomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func main() {
	logger := logging.GetLogger()

	r := gin.Default()
	r.GET("/", HomePage)

	logger.Info("run localhost:8080")
	r.Run()
}




