package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"scheduleService/model"
	"strconv"
	"sync"
)

var syncTasks = sync.Map{}

// ScheduleStart 啟動排程任務
func ScheduleStart() {
	log.Info("[ScheduleService] is running")

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
		log.Info("[ScheduleService] check job start")
		query, err := JobService().Query("")
		if err != nil {
			panic("mongo query job fail")
		}

		//for _, v := range query {
		//	//_, exists := tasks[v.ID.Hex()]
		//	idKey := strconv.FormatInt(v.ID, 10)
		//
		//	_, exists := tasks[idKey]
		//
		//	if v.Status == "running" {
		//		if !exists {
		//			fmt.Println("new task:  ", v.Name)
		//			tasks[idKey] = addTask(v)
		//			syncTasks.Store(idKey, addTask(v))
		//		}
		//	} else {
		//		if exists {
		//			fmt.Println("task: ", v.Name, " - stop")
		//			tasks[idKey].Stop()
		//			delete(tasks, idKey)
		//			syncTasks.Delete(idKey)
		//		}
		//	}
		//}

		for _, v := range query {
			idKey := FormatTaskId(v.ID, v.Name)
			task, exists := syncTasks.Load(idKey)

			if v.Status == "running" {
				if !exists {
					fmt.Println("new task:  ", v.Name)
					syncTasks.Store(idKey, addTask(v))
				}
			} else {
				if exists {
					fmt.Println("task: ", v.Name, " - stop")

					task.(*cron.Cron).Stop()
					syncTasks.Delete(idKey)
				}
			}
		}

		log.Info("[ScheduleService] check job end")
	})

	c.Start()
}

// addTask 根據 list 建立任務並執行第一次
func addTask(task *model.Job) (c *cron.Cron) {
	c = cron.New()
	f := func() {
		//logrus.Info("[ScheduleService]",task.Name, task.Path, common.MillisecondTimestamp())
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
				fmt.Println(err)
			}
			defer rsp.Body.Close()

			body, err := ioutil.ReadAll(rsp.Body)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println("RSP:", string(body))

			break
		}

	}
	f()
	c.AddFunc(task.Cron, f)
	c.Start()

	return c
}

func StopTask(id string, name string) {
	id = FormatTaskId(id, name)
	task, exists := syncTasks.Load(id)
	if exists {
		task.(*cron.Cron).Stop()
		syncTasks.Delete(id)
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
