package main

import (
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/logging"
	"github.com/gin-gonic/gin"
	Bitmex "github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/bitmex"
)

func HomePage(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World",
	})
}

func main() {
	logger := logging.GetLogger()

	r := gin.Default()
	r.LoadHTMLFiles("index.html")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.GET("/ws", func(c *gin.Context) {
		Bitmex.Wshandler(c.Writer, c.Request)
	})

	logger.Info("run localhost:8080")
	r.Run()
}




