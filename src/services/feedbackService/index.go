package feedbackService

import (
	"errors"
	"fmt"
	"os"
	cfg "todoList/config"
	"todoList/src/models/feedbackModel"
	"todoList/src/models/user"
	"todoList/src/services/common"
	"todoList/src/utils/database"
)

type FeedbackService struct{}

type CreateForm struct {
	Content string `form:"content"`
	UserId  string `form:"user_id"`
	Imgs    string `form:"imgs"`
}

func NewModel() *feedbackModel.FeedbackModel {
	return new(feedbackModel.FeedbackModel)
}
func (FeedbackService) Create(form *CreateForm) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	user := user.User()
	if user.Id == "" {
		error = errors.New("用户登陆信息异常")
		return
	}

	feedback := NewModel()
	feedback.Id = common.GetUID()
	feedback.UserId = user.Id
	feedback.Content = form.Content
	feedback.Imgs = form.Imgs

	error = db.Model(feedback).Create(feedback).Error

	return
}

func (FeedbackService) GenerateSavePath(prefix string) string {
	config, error := cfg.Config()
	if error != nil {
		return ""
	}

	savePath := config.Section("environment").Key("feedback_img_path").String()
	if savePath == "" {
		return ""
	}

	path := fmt.Sprintf("%s/%s", savePath, prefix)

	_, error = os.Stat(path)
	if error != nil {
		if os.IsNotExist(error) {
			error := os.MkdirAll(path, os.ModePerm)
			if error != nil {
				return ""
			}
		} else {
			return ""
		}
	}

	return path
}