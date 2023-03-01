package main

import (
	controller_server "simple-hosting/compositor/app/controllers/server"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	serGroup := r.Group("/server")
	{
		serGroup.GET("/list", controller_server.GetServersList)
		serGroup.POST("/create", controller_server.CreateServer)
	}
}
