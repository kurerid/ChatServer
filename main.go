package main

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}
var c chan []byte

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
		c <- message
	}

}

func receive(ws *websocket.Conn) {
	for {
		message := <-c
		fmt.Println(message)
		for currConn := range users {
			if currConn != ws {
				err := currConn.WriteMessage(websocket.TextMessage, message)
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
	go receive(conn)
}

func main() {
	c = make(chan []byte)
	http.HandleFunc("/ws", newConnection)
	log.Fatal(http.ListenAndServe(":8000", nil))
}
