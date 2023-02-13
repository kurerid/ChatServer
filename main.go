package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	c        chan *Message
	rooms    []Room
	users    map[*websocket.Conn]bool
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type Client struct {
	ws   *websocket.Conn
	room int
}

type GetRoomOutput struct {
	rooms []Room
}

type Room struct {
	id    int
	users []*websocket.Conn
}

type Message struct {
	Data []byte
	Conn *websocket.Conn
}

func main() {
	rooms = make([]Room, 10)
	rooms[0] = Room{
		id:    0,
		users: make([]*websocket.Conn, 10),
	}
	c = make(chan *Message)
	go send()
	http.HandleFunc("/ws", newConnection)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func read(cl *Client) {
	for {
		mt, message, err := cl.ws.ReadMessage()
		if err != nil {
			log.Println(err, " ЧТЕНИИЕ")
			break
		}
		if mt == websocket.CloseMessage {
			delete(users, cl.ws)
			_ = cl.ws.Close()
			return
		}
		if message[0] == '0' {
			getMenu(cl)
		}
		c <- &Message{
			Conn: cl.ws,
			Data: message,
		}
	}
}

func send() {
	for {
		message := <-c
		for currConn := range users {
			if currConn != message.Conn {
				err := currConn.WriteMessage(websocket.TextMessage, message.Data)
				if err != nil {
					log.Println(err, "ОТПРАВКА")
				}
			}
		}
	}
}

func newConnection(w http.ResponseWriter, r *http.Request) {
	if users == nil {
		users = make(map[*websocket.Conn]bool)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err, " АПГРЕЙД")
		return
	}
	cl := Client{
		ws:   conn,
		room: 0,
	}
	getMenu(&cl)
	users[conn] = true
	log.Println("New connection was upgraded")
	go read(&cl)

}

const (
	global     = "1"
	createRoom = "2"
	selectRoom = "3"
)

func getMenu(client *Client) {
	if err := client.ws.WriteMessage(websocket.TextMessage, []byte("1. Общий чат\n2. Создать комнату\n3. Выбор комнаты")); err != nil {
		log.Println(err)
	}
	_, bytes, err := client.ws.ReadMessage()
	if err != nil {
		log.Println(err)
	}
	switch string(bytes) {
	case global:
		rooms[0].users = append(rooms[0].users, client.ws)
	case createRoom:
		room := Room{
			id:    len(rooms) + 1,
			users: make([]*websocket.Conn, 10),
		}
		rooms = append(rooms, room)
		room.users = append(room.users, client.ws)
	case selectRoom:
		output := GetRoomOutput{}
		for _, room := range rooms {
			if room.id == 0 {
				continue
			}
			output.rooms = append(output.rooms, room)
		}
		bytes, err := json.Marshal(&output)
		if err != nil {
			log.Println(err)
		}
		err = client.ws.WriteMessage(websocket.TextMessage, bytes)
		if err != nil {
			log.Println(err)
		}
	default:

	}
}
