package bitmex

import (
	"encoding/json"
	"errors"
	Auth "github.com/SvyatobatkoVlad/Websocket-API-Bitmex/internal/auth"
	"github.com/SvyatobatkoVlad/Websocket-API-Bitmex/pkg/logging"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/websocket"
	"net/http/httputil"
	"net/url"
	"time"
)


var (
	// ErrCommandInvalid if command is invalid
	ErrCommandInvalid = errors.New("error provided command is invalid for Bitmex")
	// ErrSetConnection if failed to set connection
	ErrSetConnection = errors.New("error on set connection with Bitmex service")
	// ErrSendMessage if message was not sent
	ErrSendMessage = errors.New("error on send message to Bitmex")
	// ErrReadMessage if message can not be read
	ErrReadMessage = errors.New("error on read message from Bitmex")
	// ErrInvalidResponse if response is invalid
	ErrInvalidResponse = errors.New("error response is invalid from Bitmex")
)

type (
	// WebsocketClient is for websocket connection
	WebsocketClient struct {
		wsConn       *websocket.Conn
		urlToConnect string
		logger logging.Logger
	}

	// Command is for sending commands to websocket server
	Command struct {
		Op   string   `json:"op" validate:"required"`
		Args []string `json:"args,omitempty"`
	}

	// ResponseMessage is for receiving messages from websocket server
	ResponseMessage struct {
		Data      []Data `json:"data,omitempty"`
	}

	// Data is for nested values
	Data struct {
		Timestamp time.Time   `json:"timestamp,omitempty"`
		Symbol    string      `json:"symbol,omitempty"`
		LastPrice json.Number `json:"lastprice,omitempty"`
	}
)

// NewWebsocketClient to initialize WebsocketClient
func NewWebsocketClient(wsConn *websocket.Conn, urlToConnect url.URL, logger logging.Logger) *WebsocketClient {
	return &WebsocketClient{
		wsConn:       wsConn,
		urlToConnect: urlToConnect.String(),
		logger: logger,
	}
}

// SetConnection to create a connection with websocket server
func (w *WebsocketClient) SetConnection() (*WebsocketClient, error) {
	if w.wsConn != nil {
		return w, nil
	}

	conn, resp, err := websocket.DefaultDialer.Dial(w.urlToConnect, nil)
	if err != nil {
		w.logger.Warning(ErrSetConnection, ":", err)
		return nil, ErrSetConnection
	}

	w.wsConn = conn

	dumpResp, err := httputil.DumpResponse(resp, true)
	if err != nil {
		w.logger.Info("error parse Dump response")
	}

	w.logger.Info(string(dumpResp))

	////@ToDo replace this for env variables
	API_KEY := "ORqVaoVf1TJrVnKexpWjHfjk"
	API_SECRET := "mvK7p-zYF5He2eistXxXUvASoJWRGvp6eOO5TF2gn4BHI2iB"

	cmd, err := Auth.WebsocketAuthCommand(API_KEY,API_SECRET)
	if err != nil {
		w.logger.Fatalf("auth connect not work! c: %s", err.Error())
	}

	err = w.wsConn.WriteJSON(cmd)
	if err != nil {
		w.logger.Warning("error on write JSON to websocket external server for Auth: ", err)
	}

	return w, nil
}

// SendCommand to send command to the websocket server
func (w *WebsocketClient) SendCommand(message Command) error {
	validate := validator.New()
	err := validate.Struct(message)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			w.logger.Warning(e)
		}
		return ErrCommandInvalid
	}

	err = w.wsConn.WriteJSON(message)
	if err != nil {
		w.logger.Warning("error on write JSON to websocket external server: ", err)
		return ErrSendMessage
	}

	return nil
}

// ReadMessage to receive message from websocket server
func (w *WebsocketClient) ReadMessage() (*ResponseMessage, error) {
	_, msg, err := w.wsConn.ReadMessage()
	if err != nil {
		w.logger.Info(ErrReadMessage, " : ", err)
		return nil, ErrReadMessage
	}

	w.logger.Infof("got a message from Bitmex server to our app %s", msg)

	responseMessage := new(ResponseMessage)
	err = json.Unmarshal(msg, responseMessage)
	if err != nil {
		w.logger.Info(ErrInvalidResponse, " : ", msg)
		return nil, ErrInvalidResponse
	}

	return responseMessage, nil
}
