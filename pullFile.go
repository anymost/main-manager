package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gpmgo/gopm/modules/log"
	"io"
	"net/http"
	"os"
)

func PullFile(context *gin.Context) {
	var page PageConfig
	if err := context.ShouldBindJSON(&page); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "args error"})
		return
	}
	if page.Path == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "path arg cannot be empty"})
		return
	}
	if page.Url == "" {
		context.JSON(http.StatusBadRequest, gin.H{"message": "url arg cannot be empty"})
		return
	}
	resp, err := http.Get(page.Url)
	if err != nil {
		log.Error("%s", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
		return
	}
	if resp.StatusCode != http.StatusOK {
		context.JSON(http.StatusInternalServerError, gin.H{"message": resp.Status})
		return
	}
	filePath := "./template/" + page.Path + ".html"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()
	if err != nil {
		log.Error("%s", err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		io.Copy(file, resp.Body)
		redisClient.Set(page.Path, filePath, RedisDuration)
		context.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}
