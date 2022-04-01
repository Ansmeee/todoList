package main

import (
	"fmt"
	"os"
	"time"
	"todoList/config"
	"todoList/src/models/msgModel"
	"todoList/src/models/todo"
	"todoList/src/services/common"
	"todoList/src/services/msgService"
	"todoList/src/utils/database"
)

func main() {
	fmt.Println("task check command running")
	config.InitConfig()

	db := database.Connect("")
	defer database.Close(db)

	d, _ := time.ParseDuration("24h")
	sDate := time.Now().Format("2006-01-02")
	eDate := time.Now().Add(d * 3).Format("2006-01-02")

	data := make([]*todo.TodoModel, 0)
	if error := db.Model(&todo.TodoModel{}).Where("deadline between ? and ?", sDate, eDate).Find(&data).Error; error != nil {
		fmt.Println("task check command error", error.Error())
		os.Exit(0)
	}

	msgSvc := new(msgService.MsgService)
	for _, todo := range data {
		data := msgSvc.NewMsgModel()
		data.Id = common.GetUID()
		data.UserId = todo.UserId
		data.Status = msgModel.StatusUnread
		data.Link = fmt.Sprintf("/dir/%s?s_id=%s", todo.ListId, todo.Id)
		if todo.Deadline == sDate {
			data.Force = msgModel.Force
			data.Content = fmt.Sprintf("%s", todo.Title)
		} else {
			deadline, _ := time.Parse("2006-01-02", todo.Deadline)
			currentTime, _ := time.Parse("2006-01-02", sDate)
			subDay := int(deadline.Sub(currentTime).Hours() / 24)
			data.Content = fmt.Sprintf("%s 距离截至时间还剩 %d 天", todo.Title, subDay)
		}

		msgSvc.Create(data)
	}

	fmt.Println("task check command completed")
}
