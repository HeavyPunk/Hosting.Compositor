package main

import (
	"fmt"
	controller_server "simple-hosting/compositor/app/controllers/server"
	"simple-hosting/compositor/app/settings"
	file_settings_provider "simple-hosting/go-commons/settings/file-settings-provider"

	"github.com/gin-gonic/gin"
)

func main() {
	config := file_settings_provider.GetSetting[settings.ServiceSettings]("settings.yml")
	r := gin.Default()
	serGroup := r.Group("/server")
	{
		serGroup.GET("/list", controller_server.GetServersList)
		serGroup.POST("/create", controller_server.CreateServer)
	}
	r.Run(":" + fmt.Sprint(config.Socket.Port))
}