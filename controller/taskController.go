package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskControllerStruct struct{}
type fmtResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    interface{}
}

func TaskController() TaskControllerStruct {
	return TaskControllerStruct{}
}

func (r TaskControllerStruct) CreateTask(c *gin.Context) {
	taskName := c.PostForm("name")

	var validationField struct {
		Name string `form:"name" json:"name" xml:"name"  binding:"required"`
	}
	if err := c.ShouldBind(&validationField); err != nil {
		fmtResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	msg := taskName + " create task!"

	fmtResp{
		Code:    http.StatusOK,
		Message: msg,
	}.customResponse(c)
}

func (r TaskControllerStruct) QueryTask(c *gin.Context) {
	taskId := c.Query("id")

	// 建立驗證
	var validationField struct {
		Name string `form:"user" json:"user" xml:"user"  binding:"required"`
	}
	// 若有錯誤返回
	if err := c.ShouldBind(&validationField); err != nil {
		fmtResp{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}.customResponse(c)
		return
	}

	msg := "query " + taskId
	fmtResp{
		Code:    http.StatusOK,
		Message: msg,
		Data: gin.H{
			"test": "test",
		},
	}.customResponse(c)
}

// 統一 resp
func (r fmtResp) customResponse(c *gin.Context) {
	h := gin.H{
		"code":    r.Code,
		"message": r.Message,
	}

	//if r.Code == http.StatusOK {
	//	h["data"] = r.Data
	//}

	c.JSON(r.Code, h)
}
