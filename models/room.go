package models

import "github.com/gorilla/websocket"

type GetRoomOutput struct {
	Rooms []Room
}

type Room struct {
	Id    int
	Users []*websocket.Conn
}
