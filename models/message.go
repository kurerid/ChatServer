package models

import "github.com/gorilla/websocket"

type Message struct {
	Data []byte
	Conn *websocket.Conn
}
