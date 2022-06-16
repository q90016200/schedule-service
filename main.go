package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"net/http"
	"scheduleService/config"
	"time"
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

func testCron()  {

	i := 0
	c := cron.New(cron.WithSeconds())
	c1 := &cron.Cron{}
	c.AddFunc("*/1 * * * * *", func() {
		fmt.Println("c Every 1 sec", i)

		if i == 0 {
			c1 = subCron()
		}

		if i == 10 {
			fmt.Println("c1 stop")
			c1.Stop()
		}

		i++
	})
	c.Start()

	//i2 := 0
	//b := cron.New(cron.WithSeconds())
	//b.AddFunc("*/1 * * * * *", func() {
	//	fmt.Println("b Every 1s", i)
	//	i2 ++
	//	if i2 > 5 { b.Stop() }
	//})
	//b.Start()
}

func subCron() (c *cron.Cron) {
	fmt.Println("subCron")
	c = cron.New(cron.WithSeconds())
	_, err := c.AddFunc("*/1 * * * * *", func() {
		fmt.Println("subCron:", time.Now())
	})
	if err != nil {
		fmt.Println(err)
	}
	c.Start()

	return c
}