package main

import (
	"fmt"
	"os"
	controller_server "simple-hosting/compositor/app/controllers/server"
	middleware_auth "simple-hosting/compositor/app/middlewares/auth"
	middleware_logging "simple-hosting/compositor/app/middlewares/logging"
	"simple-hosting/compositor/app/settings"

	"github.com/gin-gonic/gin"
)

func main() {
	config := settings.GetServiceSettings()
	configuration := config.App.Configuration

	logFile, _ := os.Create("compositor.log")
	gin.DefaultWriter = logFile
	gin.SetMode(configuration)
	r := gin.Default()

	r.Use(middleware_auth.CheckAuth)
	r.Use(middleware_logging.ErrorLogger)

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
