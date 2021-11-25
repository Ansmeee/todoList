package list

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todoList/src/services/list"
	response "todoList/src/utils/response"
	lv "todoList/src/utils/validator/list"
)

type ListController struct {}
var service = new(list.ListService)

func (ListController) Create(request *gin.Context) {
	response := response.Response{request}
	var error error

	list := service.NewModel()
	
	error = request.ShouldBind(list)
	if error != nil {
		response.ErrorWithMSG("创建失败了，再试一次吧")
		return
	}

	validator := new(lv.ListValidator)
	error = validator.Validate(*list, lv.CreateRules)
	if error != nil {
		response.ErrorWithMSG(fmt.Sprintf("创建失败了，%s", error.Error()))
		return
	}

	data, error := service.Create(list)
	if error != nil {
		response.ErrorWithMSG("创建失败了，再试一次吧")
		return
	}

	response.SuccessWithData(data)
}