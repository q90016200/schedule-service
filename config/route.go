package config

import (
	"scheduleService/controller"

	"github.com/gin-gonic/gin"
)

func RouteJob(r *gin.Engine) {
	job := r.Group("/job")
	{
		job.POST("/", controller.JobController().Create)
		job.GET("/", controller.JobController().Query)
	}
}
