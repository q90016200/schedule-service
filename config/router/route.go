package router

import (
	"github.com/gin-gonic/gin"
	"scheduleService/controller"
)

func RouteJobs(r *gin.Engine) {
	jobs := r.Group("/api")
	{
		jobs.POST("/jobs", controller.JobController().Create)
		jobs.GET("/jobs", controller.JobController().Query)
		jobs.GET("/jobs:id", controller.JobController().Query)
		jobs.PUT("/jobs/:id", controller.JobController().Update)
		jobs.DELETE("/jobs/:id", controller.JobController().Delete)
		jobs.POST("/jobs/trigger", controller.JobController().Trigger)
	}
}
