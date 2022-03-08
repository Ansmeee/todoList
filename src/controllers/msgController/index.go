package msgController

import (
	"github.com/gin-gonic/gin"
	"todoList/src/models/msgModel"
	"todoList/src/models/user"
	"todoList/src/services/msgService"
	"todoList/src/utils/response"
)

type MsgController struct{}

var service = &msgService.MsgService{}
func (MsgController) List(request *gin.Context) {
	response := response.Response{request}

	if !user.Active() {
		response.ErrorWithMSG("请登录后再试")
		return
	}

	form := new(msgService.ListForm)
	error := request.ShouldBind(form)
	messages, error := service.List(form)
	if error != nil {
		response.ErrorWithMSG("消息加载失败")
		return
	}

	if len(messages) > 0 {
		response.SuccessWithData(messages)
		return
	}

	response.SuccessWithData([]msgModel.MsgModel{})
	return
}

func (MsgController) Unread(request *gin.Context)  {
	response := response.Response{request}

	count := service.UnreadCount()
	data := map[string]int64{"count": count}
	response.SuccessWithData(data)
}

func (MsgController) UpdateAttr(request *gin.Context) {
	response := response.Response{request}

	form := new(msgService.AttrForm)
	error := request.ShouldBindUri(form)
	if error != nil {
		response.Success()
		return
	}

	error = request.ShouldBind(form)
	if error != nil {
		response.Success()
		return
	}

	msg, error := service.FindByID(form.Id)
	if error != nil {
		response.Success()
		return
	}

	if msg.Id == "" || msg.UserId != user.User().Id {
		response.Success()
		return
	}

	if msg.Status == msgModel.StatusRead {
		response.Success()
		return
	}

	if error = service.Update(msg, form.Name, form.Value); error != nil {
		response.Success()
		return
	}

	response.Success()
	return
}
