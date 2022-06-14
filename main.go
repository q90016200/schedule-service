package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"scheduleService/config"
)

func main() {
	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// 設定靜態檔案路由
	r.Static("/public", "public")
	r.Static("/assets", "public/assets")

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Hello Gin",
		})
	})

	r.GET("/test", func(c *gin.Context) {
		c.HTML(http.StatusOK, "test.html", gin.H{
			"title": "Hello Gin",
		})
	})
	config.RouteTask(r)

	// By default it serves on :8080 unless a PORT environment variable was defined.
	err := r.Run()
	// router.Run(":3000") for a hard coded port
	if err != nil {
		panic(err)
	}
}