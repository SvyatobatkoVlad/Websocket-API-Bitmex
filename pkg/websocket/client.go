package websocket

import (
	"encoding/json"
	Bitmex "github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/bitmex"
	"github.com/gorilla/websocket"
)


type Client struct {
	ID           string
	Conn         *websocket.Conn
	WsServer         *WsServer
	BitmexClient *Bitmex.WebsocketClient
	Subscription map[string]struct{}
	send     chan []byte
}

const (
	sub = "subscribe"
	unsub = "subscribe"
)


func (c *Client) Read() {
	defer func() {
		c.WsServer.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			c.WsServer.logger.Warning(err)
			return
		}

		command := new(Bitmex.Commands)
		err = json.Unmarshal(p, command)
		if err != nil {
			c.WsServer.logger.Info("Invalid Message")
		}

		c.WsServer.logger.Infof("Got message from users to our app : %s", command)

		if command.Action == sub {
			if len(command.Symbols) == 0 {
				c.Subscription["ALL"] = struct{}{}
			} else {
				for _, symbol := range command.Symbols {
					c.Subscription[symbol] = struct{}{}
				}
			}

			_ = Bitmex.CommandExecution(c.BitmexClient, command)
		}

		if command.Action == unsub {
			if len(command.Symbols) == 0 {
				for key := range c.Subscription {
					delete(c.Subscription, key)
				}
			} else {
				for _, symbol := range command.Symbols {
					delete(c.Subscription, symbol)
				}
			}
		}
	}
}