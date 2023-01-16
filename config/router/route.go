package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"scheduleService/controller"
)

func RouteJobs(r *gin.Engine, db *gorm.DB) {
	jobs := r.Group("/api")
	{
		jobs.POST("/jobs", controller.JobController(db).Create)
		jobs.GET("/jobs", controller.JobController(db).Query)
		jobs.GET("/jobs:id", controller.JobController(db).Query)
		jobs.PUT("/jobs/:id", controller.JobController(db).Update)
		jobs.DELETE("/jobs/:id", controller.JobController(db).Delete)
		jobs.POST("/jobs/trigger", controller.JobController(db).Trigger)
	}
}
