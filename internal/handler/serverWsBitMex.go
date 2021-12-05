package handler

import (
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/bitmex"
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/logging"
	websocketSruct "github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool{
		return true
	},
}

func ServeWs(wsServer *websocketSruct.WsServer, bitmexClient *bitmex.WebsocketClient, c *gin.Context) {
	logger := logging.GetLogger()


	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Fatalf("Error updated http to websocket: %s", err)
	}

	client := &websocketSruct.Client{
		Conn:         conn,
		WsServer:         wsServer,
		BitmexClient: bitmexClient,
		Subscription: make(map[string]struct{}),
	}

	wsServer.Register <- client
	client.Read()
}