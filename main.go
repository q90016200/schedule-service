package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"net/http"
	"scheduleService/config/router"
	"scheduleService/service"
	"time"
)

func main() {
	//testCron()

	// schedule service start
	service.ScheduleStart()

	// -----------------------------

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	r := gin.Default()
	r.Use(CORSMiddleware())
	//r.Use(cors.Default())

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
	router.RouteJobs(r)

	//r.Group("/api/job")

	//r.POST("/api/jobs", controller.JobController().Create)
	//r.GET("/api/jobs", controller.JobController().Query)
	//r.GET("/api/jobs/:id", controller.JobController().Query)
	//r.PUT("/api/jobs/:id", controller.JobController().Update)
	//r.DELETE("/api/jobs/:id", controller.JobController().Delete)

	// By default it serves on :8080 unless a PORT environment variable was defined.
	r.Run()
}

func testCron() {
	i := 0
	c := cron.New(cron.WithSeconds())
	//sc := &cron.Cron{}
	tasks := make(map[string]*cron.Cron)
	c.AddFunc("*/1 * * * * *", func() {
		fmt.Println("c Every 1 sec", i)

		if i == 0 {
			tasks["1"] = subCron()
		}

		if i == 10 {
			fmt.Println("c1 stop")
			tasks["1"].Stop()
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

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
