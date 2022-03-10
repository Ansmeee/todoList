package feedbackController

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"path"
	"strings"
	"todoList/src/models/user"
	"todoList/src/services/feedbackService"
	"todoList/src/utils/response"
)

type FeedbackController struct {}

var service = &feedbackService.FeedbackService{}
func (FeedbackController) Create(request *gin.Context)  {
	response := response.Response{request}

	form := new(feedbackService.CreateForm)
	error := request.ShouldBind(form)

	if error != nil {
		response.ErrorWithMSG("吐槽失败，再来一次")
		return
	}

	multipartForm, error := request.MultipartForm()
	if error != nil {
		response.ErrorWithMSG("吐槽失败，再来一次")
		return
	}

	var imgList = []string{}
	files := multipartForm.File["imgs[]"]
	if len(files) > 0 {
		savePath := service.GenerateSavePath(user.User().Id)
		if savePath == ""{
			response.ErrorWithMSG("吐槽失败，再来一次")
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
	error = service.Create(form)
	if error != nil {
		response.ErrorWithMSG("吐槽失败，再来一次")
		return
	}

	response.Success()
}