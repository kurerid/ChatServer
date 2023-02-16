package models

import "github.com/gorilla/websocket"

type Client struct {
	Ws   *websocket.Conn
	Room int
}
