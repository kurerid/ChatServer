package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Data []byte
	Conn *websocket.Conn
}

var c chan *Message

func read(ws *websocket.Conn) {
	for {
		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
		}
		if mt == websocket.CloseMessage {
			delete(users, ws)
			_ = ws.Close()
			return
		}
		c <- &Message{
			Conn: ws,
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
					fmt.Println(err)
				}
			}
		}
	}
}

var users map[*websocket.Conn]bool

func newConnection(w http.ResponseWriter, r *http.Request) {
	if users == nil {
		users = make(map[*websocket.Conn]bool)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	users[conn] = true
	log.Println("New connection was upgraded")
	go read(conn)

}

func main() {
	c = make(chan *Message)
	go send()
	http.HandleFunc("/ws", newConnection)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
