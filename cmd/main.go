package main

import (
	Internal "github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal"
	Bitmex "github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/bitmex"
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/handler"
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/logging"
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/websocket"
	"net/url"
)

func main() {
	logger := logging.GetLogger()

	urlToBitmexWebsocket := url.URL{
		Scheme: "wss",
		Host:   "ws.testnet.bitmex.com",
		Path:   "/realtime",
	}

	//initialize WebsocketClient
	websocketClient := Bitmex.NewWebsocketClient(nil, urlToBitmexWebsocket, logger)

	websocketClient, err := websocketClient.SetConnection()
	if err != nil {
		logger.Fatal("error websocket connection was not set")
	}

	wsServer := websocket.NewWebsocketServer()
	go wsServer.Run()

	go websocket.ListenBitmex(wsServer, websocketClient)

	handler := new(handler.Handler)
	srv := new(Internal.Server)
	if err := srv.Run("8000", handler.InitRoutes(wsServer, websocketClient)); err != nil {
		logger.Fatalf("error ocured while running http server: %s", err.Error())
	}
}




