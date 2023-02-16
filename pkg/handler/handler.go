package handler

import (
	"github.com/gin-gonic/gin"
)

func InitRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/sign-up")
	router.GET("/sign-in")

	rooms := router.Group("/room")
	{
		//получить все комнаты
		rooms.GET("/list")
		//создать новую комнату
		rooms.POST("")
		//получить комнату по id
		rooms.GET("")
	}
	router.GET("/chat/ws", newConnection)
	return router
}
