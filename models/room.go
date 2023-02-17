package models

import "github.com/gorilla/websocket"

type Room struct {
	Id    int
	Users []*websocket.Conn
}

type RoomGetListOutput struct {
	Rooms []Room
}

type RoomGetByIdInput struct {
	Id int `json:"id" binding:"required,numeric"`
}

type RoomGetByIdOutput struct {
	Room Room `json:"room"`
}

type RoomCreateOutput struct {
	Room Room `json:"room"`
}
