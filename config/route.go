package config

import (
	"scheduleService/controller"

	"github.com/gin-gonic/gin"
)

func RouteTask(r *gin.Engine) {
	task := r.Group("/task")
	{
		task.POST("/", controller.TaskController().CreateTask)
		task.GET("/", controller.TaskController().QueryTask)
	}
}
