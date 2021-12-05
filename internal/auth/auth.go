package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

type Command struct {
Op   string   `json:"op" validate:"required"`
Args []interface{} `json:"args,omitempty"`
}

type sigParams struct {
	secret, method, path, body string
	expires                    time.Time
}

func calculateSignature(params *sigParams) (string, error) {
	raw := fmt.Sprintf("%s%s%d%s", params.method, params.path, params.expires.Unix(), params.body)
	sig := hmac.New(sha256.New, []byte(params.secret))

	if _, err := sig.Write([]byte(raw)); err != nil {
		return "", err
	}
	return hex.EncodeToString(sig.Sum(nil)), nil
}

func WebsocketAuthCommand(key, secret string) (*Command, error) {
	req := &sigParams{
		method:  "GET",
		path:    "/realtime",
		secret:  secret,
		body:    "",
		expires: expiryTime(),
	}
	sig, err := calculateSignature(req)
	if err != nil {
		return nil, err
	}

	//fmt.Println(" !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!  sig", sig, "expires: ", req.expires.Unix() )
	//{"op": "authKeyExpires", "args": ["ORqVaoVf1TJrVnKexpWjHfjk", 1638741658, "2333c8d56f1db38cfa6e5e3cd2992b78c96dd6dc07f80e549764a26342fa867c"]}
	//{"action": "subscribe", "symbols": ["XBTUSDT"]}

	cmd := &Command{
		Op:   "authKeyExpires",
		Args: []interface{}{key, req.expires.Unix(), sig},
	}
	return cmd, nil
}

func expiryTime() time.Time {
	return time.Now().Add(5 * time.Minute)
}