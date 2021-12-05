package websocket

import "github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/bitmex"

func ListenBitmex(wsServer *WsServer, w *bitmex.WebsocketClient) {
	for {
		response, err := w.ReadMessage()
		if err != nil {
			continue
		}

		wsServer.BitmexBroadcast <- response
	}

}
