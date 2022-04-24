package feedbackService

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"os"
	"path"
	"strings"
	cfg "todoList/config"
	"todoList/src/models/feedbackModel"
	"todoList/src/models/user"
	"todoList/src/services/common"
	"todoList/src/utils/database"
	"todoList/src/utils/redis"
)

type FeedbackService struct{}

type CreateForm struct {
	Content string `form:"content"`
	UserId  string `form:"user_id"`
	Imgs    string `form:"imgs"`
	Files   []*multipart.FileHeader
}

var ctx = context.Background()

func NewModel() *feedbackModel.FeedbackModel {
	return new(feedbackModel.FeedbackModel)
}

func (*FeedbackService) FeedbackFrequently() bool {
	client := redis.Connect()
	defer redis.Close(client)

	user := user.User()
	if user.Id == "" {
		return false
	}

	key := fmt.Sprintf("feedback:num:%s", user.Id)

	num, err := client.Get(ctx, key).Int()
	if err != nil {
		return false
	}

	return num >= 5
}

func (FeedbackService) Create(form *CreateForm, request *gin.Context) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	client := redis.Connect()
	defer redis.Close(client)

	user := user.User()
	if user.Id == "" {
		error = errors.New("用户登陆信息异常")
		return
	}

	files := form.Files
	var imgList = []string{}
	if len(files) > 0 {
		savePath := generateSavePath(user.Id)
		if savePath == ""{
			error = errors.New("图片保存失败")
			return
		}

		for _, file := range files {
			fileName := fmt.Sprintf("%s%s", fmt.Sprintf("%x", md5.Sum([]byte(path.Base(file.Filename)))), path.Ext(file.Filename))
			filePath := fmt.Sprintf("%s/%s", savePath, fileName)
			request.SaveUploadedFile(file, filePath)
			imgList = append(imgList, fileName)
		}
	}

	form.Imgs = strings.Join(imgList, ";")
	key := fmt.Sprintf("feedback:num:%s", user.Id)

	if err := client.Incr(ctx, key).Err(); err != nil {
		error = errors.New("提交失败，请再试一次")
		return
	}

	feedback := NewModel()
	feedback.Id = common.GetUID()
	feedback.UserId = user.Id
	feedback.Content = form.Content
	feedback.Imgs = form.Imgs

	return db.Model(feedback).Create(feedback).Error
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
