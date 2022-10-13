package router

import (
	"github.com/gin-gonic/gin"
	"scheduleService/controller"
)

func RouteJob(r *gin.Engine) {
	job := r.Group("/job")
	{
		job.POST("/", controller.JobController().Create)
		job.GET("/", controller.JobController().Query)
		job.GET("/:id", controller.JobController().Query)
		job.PUT("/:id", controller.JobController().Update)
		job.DELETE("/:id", controller.JobController().Delete)
	}
}
