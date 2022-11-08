package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"scheduleService/model"
	"strconv"
	"sync"
	"time"
)

var syncTasks = sync.Map{}
var taskLogger *log.Logger

func init() {
	layout := "2006-01-02"
	formatTime := time.Now().Format(layout)
	folderName := "logs"
	folderPath := folderName
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.Mkdir(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}

	taskLogFile, err := os.OpenFile("./logs/task-"+formatTime+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("開啟 task 日誌檔案失敗：", err)
	}

	taskLogger = log.New()
	taskLogger.Formatter = &log.JSONFormatter{}
	taskLogger.SetOutput(taskLogFile)
}

// ScheduleStart 啟動排程任務
func ScheduleStart() {
	//log.Info("[ScheduleService] is running")

	// 存放需執行的任務
	//tasks := make(map[string]*cron.Cron)

	//c := cron.New(cron.WithSeconds())
	//c := cron.New()

	logger := &CLog{clog: log.New()}
	logger.clog.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(logger)))

	// 建立檢查任務狀態任務
	c.AddFunc("* * * * *", func() {
		//log.Info("[ScheduleService] check job start")
		query, err := JobService().Query("")
		if err != nil {
			panic("query job fail")
		}

		for _, v := range query {
			id := strconv.FormatInt(v.ID, 10)
			idKey := FormatTaskId(id, v.Name)
			_, exists := syncTasks.Load(idKey)

			if v.Status == "running" {
				if !exists {
					CreateCronTask(id, v)
				}
			} else {
				if exists {
					StopCronTask(id, v.Name)
				}
			}
		}

		//log.Info("[ScheduleService] check job end")
	})

	c.Start()
}

// CreateCronTask 建立排程任務並執行第一次
func CreateCronTask(id string, task *model.Job) {
	taskId := FormatTaskId(id, task.Name)
	c := cron.New()
	f := func() {
		fmt.Println("new task:  ", task.Name)
		log.WithFields(log.Fields{
			"name":   task.Name,
			"method": task.Method,
			"path":   task.Path,
		}).Info()

		switch task.Method {
		case "http":
			url := task.Path
			client := http.Client{}
			rsp, err := client.Get(url)
			if err != nil {
				//fmt.Println(err)
				taskLogger.WithFields(log.Fields{
					"name":  task.Name,
					"path":  task.Path,
					"error": err.Error(),
				}).Error()
			}
			defer rsp.Body.Close()

			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				log.Error(task.Name + " | " + url + " | " + err.Error())
				taskLogger.WithFields(log.Fields{
					"name":  task.Name,
					"path":  task.Path,
					"error": err.Error(),
				}).Error()
			}
			fmt.Println(task.Name+" | "+url+" | ", string(body))
			taskLogger.WithFields(log.Fields{
				"name":     task.Name,
				"path":     task.Path,
				"response": string(body),
			}).Info()

			break
		}

	}
	f()
	c.AddFunc(task.Cron, f)
	c.Start()

	syncTasks.Store(taskId, c)
}

func StopCronTask(id string, name string) {
	id = FormatTaskId(id, name)
	task, exists := syncTasks.Load(id)
	if exists {
		task.(*cron.Cron).Stop()
		syncTasks.Delete(id)

		fmt.Println("task: ", name, " - stop")
	}
}

type CLog struct {
	clog *log.Logger
}

func (l *CLog) Info(msg string, keysAndValues ...interface{}) {
	l.clog.WithFields(log.Fields{
		"data": keysAndValues,
	}).Info(msg)
}

func (l *CLog) Error(err error, msg string, keysAndValues ...interface{}) {
	l.clog.WithFields(log.Fields{
		"msg":  msg,
		"data": keysAndValues,
	}).Warn(msg)
}

func FormatTaskId(id interface{}, name string) string {
	taskId := ""
	reg := regexp.MustCompile(`[\s\p{Zs}]{1,}`)
	name = reg.ReplaceAllString(name, "")

	switch id.(type) {
	case int64:
		taskId = strconv.FormatInt(id.(int64), 10)
	case string:
		taskId = id.(string)
	}
	return "task-" + taskId + "-" + name
}
