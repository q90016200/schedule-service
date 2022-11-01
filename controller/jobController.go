package controller

import (
	"fmt"
	"net/http"
	"scheduleService/model"
	"scheduleService/service"
	"time"

	"github.com/gin-gonic/gin"
)

type JobControllerStruct struct{}
type formatResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

func JobController() JobControllerStruct {
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
	}
	// 驗證請求資料
	if err := c.ShouldBind(&requestField); err != nil {
		respFmt.Code = http.StatusBadRequest
		respFmt.Message = err.Error()
		respFmt.customResponse(c)
		return
	}
	// 寫入任務至 database
	err := service.JobService().Create(model.Job{
		Name:   requestField.Name,
		Method: requestField.Method,
		Path:   requestField.Path,
		Cron:   requestField.Cron,
	})
	if err != nil {
		respFmt.Code = http.StatusBadRequest
		respFmt.Message = err.Error()
		respFmt.customResponse(c)
		return
	}

	respFmt.customResponse(c)
}

func (r JobControllerStruct) Query(c *gin.Context) {
	id := c.DefaultQuery("id", "")
	result, err := service.JobService().Query(id)
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

	data := model.Job{
		Name:      c.PostForm("name"),
		Method:    c.PostForm("method"),
		Path:      c.PostForm("path"),
		Cron:      c.PostForm("cron"),
		Status:    c.PostForm("status"),
		UpdatedAt: time.Now().UTC(),
	}

	fmt.Println(data.Name)

	err := service.JobService().Update(id, data)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	formatResp{
		Code: http.StatusOK,
	}.customResponse(c)
}

func (r JobControllerStruct) Delete(c *gin.Context) {
	id := c.Param("id")

	result, err := service.JobService().Query(id)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	err = service.JobService().Delete(id)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	// 關閉 cron 服務
	service.StopTask(id, result[0].Name)

	formatResp{
		Code: http.StatusOK,
	}.customResponse(c)
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
