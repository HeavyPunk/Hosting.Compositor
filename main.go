package main

import (
	"fmt"
	controller_server "simple-hosting/compositor/app/controllers/server"
	middleware_auth "simple-hosting/compositor/app/middlewares/auth"
	"simple-hosting/compositor/app/settings"

	"github.com/gin-gonic/gin"
)

func main() {
	config := settings.GetServiceSettings()
	configuration := config.App.Configuration

	gin.SetMode(configuration)
	r := gin.Default()

	r.Use(middleware_auth.CheckAuth)

	serGroup := r.Group("/server")
	{
		serGroup.GET("/list", controller_server.GetServersList)
		serGroup.POST("/create", controller_server.CreateServer)
		serGroup.POST("/start", controller_server.StartServer)
		serGroup.POST("/stop", controller_server.StopServer)
		serGroup.POST("/remove", controller_server.DeleteServer)
	}
	r.Run(":" + fmt.Sprint(config.App.Port))
}
