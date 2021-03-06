package feedbackController

import (
	"github.com/gin-gonic/gin"
	"todoList/src/services/feedbackService"
	"todoList/src/utils/response"
)

type FeedbackController struct {}

var service = &feedbackService.FeedbackService{}
func (FeedbackController) Create(request *gin.Context)  {
	response := response.Response{request}

	res := service.FeedbackFrequently()
	if res == true {
		response.ErrorWithMSG("提交太频繁了，稍后再试吧")
		return
	}

	form := new(feedbackService.CreateForm)
	error := request.ShouldBind(form)

	if error != nil {
		response.ErrorWithMSG("")
		return
	}

	multipartForm, error := request.MultipartForm()
	if error != nil {
		response.ErrorWithMSG("")
		return
	}

	form.Files = multipartForm.File["imgs[]"]
	fb, error := service.Create(form, request)
	if error != nil {
		response.ErrorWithMSG("")
		return
	}

	service.SendMSG2M(fb)
	response.Success()
}