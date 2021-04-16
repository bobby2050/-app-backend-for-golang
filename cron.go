package main

import (
	"github.com/robfig/cron"
	"log"
	"meizi/src/web/Model"
	"time"
)

func main()  {
	log.Println("计划任务开始 Starting...")
	c := cron.New()

	c.AddFunc("* */1 * * * *", func() {
		log.Println("回复审核...")
		Model.Reply()
	})
	//
	c.AddFunc("*/30 * * * * *", func() {
		log.Println("文章发布审核")
		Model.PublishC()
	})




	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}