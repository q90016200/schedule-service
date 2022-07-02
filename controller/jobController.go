package controller

import (
	"fmt"
	"net/http"
	"scheduleService/model"
	"scheduleService/service"

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
	// 建立驗證
	var requestField struct {
		Id string `form:"id" json:"id" xml:"id"  binding:"required"`
	}
	// 若有錯誤返回
	if err := c.ShouldBind(&requestField); err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	fmt.Println("id:", requestField.Id)

	result,err := service.JobService().Query(requestField.Id)
	if err != nil {
		formatResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	formatResp{
		Code:    http.StatusOK,
		Data: result,
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
