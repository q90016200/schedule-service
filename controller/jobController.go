package controller

import (
	"fmt"
	"gorm.io/gorm"
	"net/http"
	"scheduleService/model"
	"scheduleService/service"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type JobControllerStruct struct{}
type formatResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

var sqlDB *gorm.DB

func JobController(db *gorm.DB) JobControllerStruct {
	sqlDB = db
	return JobControllerStruct{}
}

// Create 建立任務
func (r JobControllerStruct) Create(c *gin.Context) {
	respFmt := formatResp{
		Code:    http.StatusOK,
		Message: "",
	}
	var requestField struct {
		Name   string `form:"name" json:"name" xml:"name"  binding:"required"`
		Method string `form:"method" json:"method" xml:"method" binding:"required"`
		Path   string `form:"path" json:"path" xml:"path"  binding:"required"`
		Cron   string `form:"cron" json:"cron" xml:"cron"  binding:"required"`
		Status string `form:"status" json:"status" xml:"status"  binding:"required"`
		Group  string `form:"group" json:"group"`
	}
	// 驗證請求資料
	if err := c.ShouldBind(&requestField); err != nil {
		respFmt.Code = http.StatusBadRequest
		respFmt.Message = err.Error()
		respFmt.customResponse(c)
		return
	}

	jobName := strings.TrimSpace(requestField.Name)
	path := strings.TrimSpace(requestField.Path)
	// 寫入任務至 database
	job := model.Job{
		Name:   jobName,
		Method: requestField.Method,
		Path:   path,
		Status: requestField.Status,
		Cron:   requestField.Cron,
	}
	if requestField.Group != "" {
		job.Group = requestField.Group
	}

	jobId, err := service.JobService(sqlDB).Create(job)
	if err != nil {
		respFmt.Code = http.StatusBadRequest
		respFmt.Message = err.Error()
		respFmt.customResponse(c)
		return
	}

	// 開啟 cron 服務
	id := strconv.FormatInt(jobId, 10)
	service.CreateCronTask(id, &job)

	respFmt.customResponse(c)
}

func (r JobControllerStruct) Query(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	result, err := service.JobService(sqlDB).Query(id)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	formatResp{
		Code: http.StatusOK,
		Data: result,
	}.customResponse(c)
}

func (r JobControllerStruct) Update(c *gin.Context) {
	id := c.Param("id")
	jobName := strings.TrimSpace(c.PostForm("name"))
	cron := c.PostForm("cron")
	status := c.PostForm("status")
	method := c.PostForm("method")
	path := c.PostForm("path")
	path = strings.TrimSpace(path)

	data := map[string]interface{}{
		"Name":      jobName,
		"Method":    method,
		"Path":      path,
		"Cron":      cron,
		"Group":     c.DefaultPostForm("group", ""),
		"Status":    c.PostForm("status"),
		"UpdatedAt": time.Now().UTC(),
	}

	fmt.Printf("%+v", data)

	query, err := service.JobService(sqlDB).Query(id)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	err = service.JobService(sqlDB).Update(id, data)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	// 判斷排程是否重啟
	cronStart := false
	cronStop := false
	job := query[0]
	if job.Cron != cron {
		cronStop = true
		cronStart = true
	}
	if status != job.Status {
		if status == "running" {
			cronStart = true
		} else if status == "stopped" {
			cronStart = false
			cronStop = true
		}
	}
	if job.Status == "stopped" {
		cronStart = false
	}

	//fmt.Println("cronStart:", cronStart, "cronStop:", cronStop)

	if cronStop {
		service.StopCronTask(id, jobName)
	}
	if cronStart {
		service.CreateCronTask(id, &model.Job{
			Name:   jobName,
			Method: method,
			Path:   path,
			Cron:   cron,
		})
	}

	formatResp{
		Code: http.StatusOK,
	}.customResponse(c)
}

func (r JobControllerStruct) Delete(c *gin.Context) {
	id := c.Param("id")

	result, err := service.JobService(sqlDB).Query(id)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	err = service.JobService(sqlDB).Delete(id)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	// 關閉 cron 服務
	service.StopCronTask(id, result[0].Name)

	formatResp{
		Code: http.StatusOK,
	}.customResponse(c)
}

func (r JobControllerStruct) Trigger(c *gin.Context) {
	respFmt := formatResp{
		Code:    http.StatusOK,
		Message: "",
	}

	var requestField struct {
		Path      string `form:"path" json:"path" xml:"path"  binding:"required"`
		Frequency string `form:"frequency" json:"frequency" xml:"frequency"  binding:"required"`
		Interval  string `form:"interval" json:"interval" xml:"interval"  binding:"required"`
	}
	// 驗證請求資料
	if err := c.ShouldBind(&requestField); err != nil {
		respFmt.Code = http.StatusBadRequest
		respFmt.Message = err.Error()
		respFmt.customResponse(c)
		return
	}

	fmt.Printf("%+v", requestField)
}

// 統一 resp
func (r formatResp) customResponse(c *gin.Context) {
	h := gin.H{
		"code":    r.Code,
		"message": r.Message,
	}

	if r.Code == http.StatusOK {
		h["data"] = r.Data
	}

	c.JSON(r.Code, h)
}
