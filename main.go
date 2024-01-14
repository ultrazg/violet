package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"violet/config"
	"violet/server"
	"violet/util"
)

func init() {
	config.ConfigInit()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(cors.Default())

	util.P(config.Config.Server.Host, config.Config.Server.Port, Version)

	router.GET("/", server.Test) // test router
	router.POST("/video2text", server.Video2Text)
	router.POST("/removeWatermark", server.RemoveWatermark)

	err := router.Run(":" + config.Config.Server.Port)
	if err != nil {
		panic("server start fail")

		return
	}
}
