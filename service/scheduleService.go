package service

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"scheduleService/model"
)

// ScheduleStart 啟動排程任務
func ScheduleStart() {
	logrus.Info("[ScheduleService] is running")

	// 存放需執行的任務
	tasks := make(map[string]*cron.Cron)
	//c := cron.New(cron.WithSeconds())
	c := cron.New()

	// 建立檢查任務狀態任務
	c.AddFunc("*/1 * * * *", func() {
		logrus.Info("[ScheduleService] check job start")
		query, err := JobService().Query("")
		if err != nil {
			panic("mongo query job fail")
		}

		for _,v := range query{
			_, exists := tasks[v.ID.Hex()]
			if v.Status == "1" {
				if !exists {
					fmt.Println("new task:", v.Name)
					tasks[v.ID.Hex()] = newTask(v)
				}
			} else {
				if exists {
					fmt.Println("task:", v.Name, "stop")
					tasks[v.ID.Hex()].Stop()
					delete(tasks, v.ID.Hex())
				}
			}
		}
		logrus.Info("[ScheduleService] check job end")
	})

	c.Start()
}

// newTask 根據 list 建立任務並執行第一次
func newTask(task *model.Job) (c *cron.Cron) {
	c = cron.New()
	f := func() {
		//logrus.Info("[ScheduleService]",task.Name, task.Path, common.MillisecondTimestamp())
		logrus.WithFields(logrus.Fields{
			"name": task.Name,
			"method": task.Method,
			"path": task.Path,
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







