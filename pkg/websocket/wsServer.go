package websocket

import (
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/bitmex"
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/logging"
)

type WsServer struct {
	Register        chan *Client
	Unregister      chan *Client
	Clients         map[*Client]bool
	Broadcast       chan bitmex.Commands
	BitmexBroadcast chan *bitmex.ResponseMessage
	logger logging.Logger
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		Register:        make(chan *Client),
		Unregister:      make(chan *Client),
		Clients:         make(map[*Client]bool),
		Broadcast:       make(chan bitmex.Commands),
		BitmexBroadcast: make(chan *bitmex.ResponseMessage),
		logger: logging.GetLogger(),
	}
}

func (wsServer *WsServer) Run() {
	for {
		select {
		case client := <-wsServer.Register:
			wsServer.Clients[client] = true
			wsServer.logger.Info("Size of Connection WsServer: ", len(wsServer.Clients))
			for client, _ := range wsServer.Clients {
				wsServer.logger.Info("Connected client: ", client)
			}
		case client := <-wsServer.Unregister:
			delete(wsServer.Clients, client)
			wsServer.logger.Info("Size of Connection WsServer: ", len(wsServer.Clients))
			for client, _ := range wsServer.Clients {
				wsServer.logger.Info("Disconnected client: ", client)
			}
		case message := <-wsServer.Broadcast:
			wsServer.logger.Info("Sending message to all clients in wsServer")
			for client, _ := range wsServer.Clients {
				wsServer.logger.Info("Client: ", client, "Message :", message)
			}
		case Message := <-wsServer.BitmexBroadcast:
			wsServer.logger.Info("Sending message to clients in WsServer")
			for client, _ := range wsServer.Clients {
				for _, data := range Message.Data {
					_, all := client.Subscription["ALL"]
					if _, ok := client.Subscription[data.Symbol]; (ok || all) && len(data.LastPrice) != 0 {
						if err := client.Conn.WriteJSON(data); err != nil {
							wsServer.logger.Warning(err)
							return
						}
					}
				}
			}
		}
	}
}