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
		recvMSG()
		time.Sleep(5 * time.Second)
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
	content := fmt.Sprintf("用户【%s】提交了新的用户反馈，请及时处理", user.Email)
	mSVC := new(mailSVC.MailSVC)
	if err := mSVC.SendText(subject, content, receivers...); err != nil {
		log.Println("feedback watcher error:", err)
		return
	}

	return
}