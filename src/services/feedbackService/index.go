package feedbackService

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path"
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
func (FeedbackService) Create(request *gin.Context, form *CreateForm, fileList []*multipart.FileHeader) (error error) {
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


	var imgList = []string{}
	if len(fileList) > 0 {
		savePath := generateSavePath(user.Id)
		if savePath == ""{
			error = errors.New("反馈文件系统异常")
			return
		}


		for _, file := range fileList {
			fileName := fmt.Sprintf("%s%s", fmt.Sprintf("%x", md5.Sum([]byte(path.Base(file.Filename)))), path.Ext(file.Filename))
			filePath := fmt.Sprintf("%s/%s", fileName)
			error = request.SaveUploadedFile(file, filePath)
			fmt.Println(error.Error())
		}
	}

	feedback.Imgs
	error = db.Model(feedback).Create(feedback).Error
	return
}

func generateSavePath(prefix string) string {
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