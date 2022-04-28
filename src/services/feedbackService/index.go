package feedbackService

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"
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

func (*FeedbackService) FindByID(id string) (*feedbackModel.FeedbackModel, error) {
	db := database.Connect("")
	defer database.Close(db)

	fb := NewModel()
	err := db.Where("uid = ?", id).Find(fb).Error
	return fb, err
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

func (FeedbackService) SendMSG2M(fb *feedbackModel.FeedbackModel) {
	if fb.Id == "" {
		return
	}

	client := redis.Connect()
	defer redis.Close(client)

	client.LPush(ctx, "feedback:msg:list", fb.Id)
}

func (FeedbackService) Create(form *CreateForm, request *gin.Context) (*feedbackModel.FeedbackModel, error) {
	db := database.Connect("")
	defer database.Close(db)

	client := redis.Connect()
	defer redis.Close(client)

	user := user.User()
	if user.Id == "" {
		return nil, errors.New("用户登陆信息异常")
	}

	files := form.Files
	var imgList = []string{}
	if len(files) > 0 {
		savePath := generateSavePath(user.Id)
		if savePath == ""{
			return nil, errors.New("图片保存失败")
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
		return nil, errors.New("提交失败，请再试一次")
	}

	expireAt := time.Now().Add(1 * time.Hour)
	if _, err := client.ExpireAt(ctx, key, expireAt).Result(); err != nil {
		log.Println("FeedbackService Create Error:", err)
		return nil, errors.New("提交失败，请再试一次")
	}

	feedback := NewModel()
	feedback.Id = common.GetUID()
	feedback.UserId = user.Id
	feedback.Content = form.Content
	feedback.Imgs = form.Imgs

	if err := db.Model(feedback).Create(feedback).Error; err != nil {
		log.Println("FeedbackService Create Error:", err)
		return nil, err
	}

	return feedback, nil
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
