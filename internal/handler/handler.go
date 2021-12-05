package handler

import (
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/bitmex"
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/websocket"
	"github.com/gin-gonic/gin"
)



type Handler struct {

}

func (h *Handler) InitRoutes(wsServer *websocket.WsServer, bitmexClient *bitmex.WebsocketClient) *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		bitMex := api.Group("BitMex")
		{
			bitMex.GET("/ws", func(c *gin.Context) {
				ServeWs(wsServer, bitmexClient, c)
			})
		}
	}

	return router
}
