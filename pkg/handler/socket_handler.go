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
	users    map[models.Client]bool
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func newConnection(c *gin.Context) {
	if users == nil {
		users = make(map[models.Client]bool)
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}
	cl := models.Client{
		Ws:   conn,
		Room: 0,
	}
	users[cl] = true
	for {
		_, bytes, err := cl.Ws.ReadMessage()
		if err != nil {
			if closeError, ok := err.(*websocket.CloseError); ok {
				if websocket.IsCloseError(closeError, closeError.Code) {

					break
				}
			}
			log.Println(err)
			continue
		}
		message := models.Message{
			Data: bytes,
			Conn: cl.Ws,
		}

		if err = send(&message); err != nil {
		}

	}
}

func send(message *models.Message) error {
	for currUser := range users {
		if currUser.Ws != message.Conn {
			if err := currUser.Ws.WriteMessage(websocket.TextMessage, message.Data); err != nil {
				log.Println(err)
			}
		}
	}
	return nil
}
