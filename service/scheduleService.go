package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/zeromicro/go-zero/zrpc"
	_ "github.com/zeromicro/zero-contrib/zrpc/registry/consul"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
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

	// 超過天數不保留 log
	removeDate := time.Now().AddDate(0, 0, -10).Format(layout)
	os.RemoveAll("./logs/task-" + removeDate + ".log")

	taskLogger = log.New()
	taskLogger.Formatter = &log.JSONFormatter{}
	taskLogger.SetOutput(taskLogFile)
}

// ScheduleStart 啟動排程任務
func ScheduleStart(db *gorm.DB) {
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
		query, err := JobService(db).Query("")
		if err != nil {
			panic("query job fail")
		}

		for _, v := range query {
			id := strconv.FormatInt(v.ID, 10)
			idKey := FormatTaskId(id)
			_, exists := syncTasks.Load(idKey)

			if v.Status == "running" {
				if !exists {
					//createStatus := true
					//if v.Method == "http" {
					//	u, err := url.Parse(v.Path)
					//	if err != nil {
					//		createStatus = false
					//	}
					//	fmt.Println(u.Host, u.Path)
					//}
					//if createStatus {
					CreateCronTask(id, v)
					//}
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
	taskId := FormatTaskId(id)
	fmt.Println("new task:  ", task.Name)

	//defer func() {
	//	// 可以取得 panic 的回傳值
	//	r := recover()
	//	if r != nil {
	//		fmt.Println("Recovered in f", r)
	//	}
	//}()

	logger := &CLog{clog: log.New()}
	logger.clog.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	c := cron.New(cron.WithChain(cron.SkipIfStillRunning(logger), cron.Recover(logger)))

	//c := cron.New()
	f := func() {
		//log.WithFields(log.Fields{
		//	"name":   task.Name,
		//	"method": task.Method,
		//	"group":  task.Group,
		//	"path":   task.Path,
		//}).Info()

		switch task.Method {
		case "http":
			requestUrl := task.Path
			client := http.Client{Timeout: time.Second * 20}
			rsp, err := client.Get(requestUrl)
			if err != nil {
				//fmt.Println(err)
				taskLogger.WithFields(log.Fields{
					"name":  task.Name,
					"path":  task.Path,
					"group": task.Group,
					"error": err.Error(),
				}).Error()
				return
			}
			defer rsp.Body.Close()

			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				log.Error(task.Name + " | " + requestUrl + " | " + err.Error())
				taskLogger.WithFields(log.Fields{
					"name":  task.Name,
					"path":  task.Path,
					"group": task.Group,
					"error": err.Error(),
				}).Error()
			}
			//fmt.Println(task.Name+" | "+requestUrl+" | ", string(body))
			taskLogger.WithFields(log.Fields{
				"name":     task.Name,
				"path":     task.Path,
				"group":    task.Group,
				"response": string(body),
			}).Info()

			break
		case "grpc":
			conf := zrpc.RpcClientConf{
				Target:  task.Consul,
				Timeout: 20000,
			}
			client, _ := zrpc.NewClient(conf)
			g := client.Conn()
			em := empty.Empty{}
			err := g.Invoke(context.Background(), task.Path, &em, &em)
			log.Error(task.Name + " | " + task.Consul + task.Path + " | " + err.Error())
			defer g.Close()
			//os.Exit(0)

			break
		}

	}
	//f()
	c.AddFunc(task.Cron, f)
	c.Start()

	syncTasks.Store(taskId, c)
}

func StopCronTask(id string, name string) {
	id = FormatTaskId(id)
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

func FormatTaskId(id interface{}) string {
	taskId := ""
	//reg := regexp.MustCompile(`[\s\p{Zs}]{1,}`)
	//name = reg.ReplaceAllString(name, "")

	switch id.(type) {
	case int64:
		taskId = strconv.FormatInt(id.(int64), 10)
	case string:
		taskId = id.(string)
	}
	//return "task-" + taskId + "-" + name
	return "task-" + taskId
}
