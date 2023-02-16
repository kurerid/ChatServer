package handler

import (
	"ChatServer/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var (
	c        chan *models.Message
	rooms    []models.Room
	users    map[*websocket.Conn]bool
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func newConnection(c *gin.Context) {
	if users == nil {
		users = make(map[*websocket.Conn]bool)
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err, " АПГРЕЙД")
		return
	}
	cl := models.Client{
		Ws:   conn,
		Room: 0,
	}
	users[conn] = true
	log.Println("New connection was upgraded")
	go read(&cl)

}

func read(cl *models.Client) {
	for {
		mt, message, err := cl.Ws.ReadMessage()
		if err != nil {
			log.Println(err, " ЧТЕНИИЕ")
			break
		}
		if mt == websocket.CloseMessage {
			delete(users, cl.Ws)
			_ = cl.Ws.Close()
			return
		}
		c <- &models.Message{
			Conn: cl.Ws,
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
