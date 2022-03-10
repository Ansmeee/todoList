package feedbackController

import (
	"github.com/gin-gonic/gin"
	"mime/multipart"
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

	var fileList = []*multipart.FileHeader{}
	files := multipartForm.File["imgs[]"]
	for _, file := range files{
		fileList = append(fileList, file)
	}

	error = service.Create(request, form, fileList)
	if error != nil {
		response.ErrorWithMSG("吐槽失败，再来一次")
		return
	}

	response.Success()
}