package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

const RedisDuration = time.Hour * 24 * 365 * 100

type PageConfig struct {
	Path string `json:"path"`
	Url string `json:"url"`
}

func main() {
	redisClient := StartRedis()
	router := gin.Default()
	router.POST("/pull", PullFile)

	router.GET("/:pagePath", func(context *gin.Context) {
		fmt.Println("get")
		pagePath := context.Param("pagePath")
		notFoundPath := redisClient.Get("404")
		templatePath := redisClient.Get(pagePath)
		if templatePath.Val() != "" {
			context.File(templatePath.Val())
		} else {
			context.File(notFoundPath.Val())
		}
	})
	router.Run(":80")
}
