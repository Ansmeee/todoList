package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"
	"todoList/config"
	"todoList/src/services/feedbackService"
	"todoList/src/services/mailSVC"
	"todoList/src/services/user"
	"todoList/src/utils/redis"
)

var ctx = context.Background()

func main() {
	fmt.Println("feedback watcher running...")
	config.InitConfig()

	for {
		go recvMSG()
		time.Sleep(1 * time.Second)
	}
}

func recvMSG() {
	client := redis.Connect()
	defer redis.Close(client)

	id, err := client.RPop(ctx, "feedback:msg:list").Result()
	if err != nil {
		log.Println("feedback watcher error:", err)
		return
	}

	fbSVC := new(feedbackService.FeedbackService)
	fb, err := fbSVC.FindByID(id)

	if err != nil {
		log.Println("feedback watcher error:", err)
		return
	}

	if fb.Id == "" {
		return
	}

	cfg, err := config.Config()
	if err != nil {
		log.Println("feedback watcher error:", err)
		return
	}

	recvs := cfg.Section("feedback").Key("operator").String()
	receivers := strings.Split(recvs, ";")

	userSVC := new(user.UserService)
	err, user := userSVC.FindByID(fb.UserId)
	if err != nil || user.Id == "" {
		log.Println("feedback watcher error:", err)
		return
	}

	subject := "土豆清单（ToDoo）用户反馈"
	cons := make([]string, 0)
	cons = append(cons, "用户反馈信息：")
	cons = append(cons, fmt.Sprintf("反馈用户：%s", user.Email))
	cons = append(cons, fmt.Sprintf("反馈内容：%s", fb.Content))
	cons = append(cons, fmt.Sprintf("反馈时间：%s", time.Now().Format("2006-01-02")))
	cons = append(cons, "请及时处理")
	content := strings.Join(cons, "\n")

	fmt.Println(content)
	mSVC := new(mailSVC.MailSVC)
	if err := mSVC.SendText(subject, content, receivers...); err != nil {
		log.Println("feedback watcher error:", err)
		return
	}

	return
}
