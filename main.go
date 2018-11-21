package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const RedisDuration = time.Hour * 24 * 365 * 100

type PageConfig struct {
	Path string `json:"path"`
	Url  string `json:"url"`
}

func main() {
	initRedis()
	router := gin.Default()
	router.POST("/api/files", fetchFiles)
	router.POST("/api/create", createFile)
	router.DELETE("/api/delete/:path", deleteFile)
	router.GET("/:pagePath", renderPage)
	router.Run(":80")
}

func fetchFiles(context *gin.Context) {
	pageLength := redisClient.HLen("pageList").Val()
	if pageLength > 0 {
		pageMap := redisClient.HGetAll("pageList")
		pageList := make([]*PageConfig, pageLength)
		index := 0
		for path, url := range pageMap.Val() {
			pageList[index] = &PageConfig{Path: path, Url: url}
			index++
		}
		context.JSON(http.StatusOK, gin.H{
			"data":    pageList,
			"message": "ok",
		})
	} else {
		context.JSON(http.StatusOK, gin.H{
			"data":    make([]PageConfig, 0),
			"message": "ok",
		})
	}
}

func renderPage(context *gin.Context) {
	context.Header("Server", "gin")
	pagePath := context.Param("pagePath")
	notFoundPath := redisClient.HGet("pageList", "404")
	templatePath := redisClient.HGet("pageList", pagePath)
	if templatePath.Val() != "" {
		context.File(templatePath.Val())
	} else {
		context.File(notFoundPath.Val())
	}
}

func createFile(context *gin.Context) {
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
		log.Println(err)
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
		log.Println(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		io.Copy(file, resp.Body)
		redisClient.HSet("pageList", page.Path, filePath)
		context.JSON(http.StatusOK, gin.H{"message": "ok"})
	}
}

func deleteFile(context *gin.Context) {
	path := context.Param("path")
	if path == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "path cannot be null",
		})
		return
	}
	url := redisClient.HGet("pageList", path)
	if url.Val() == "" {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "path not exist",
		})
		return
	}
	redisClient.HDel("pageList", path)
	err := os.Remove(url.Val())
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{
			"message": "delete failed",
		})
		return
	}
	context.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
}
