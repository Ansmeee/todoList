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

func (ListController) Update(request *gin.Context)  {
	response := response.Response{request}
	var error error

	data := service.NewModel()
	error = request.ShouldBind(data)
	if error != nil {
		response.ErrorWithMSG("保存失败了，再试一次吧")
		return
	}

	validator := new(lv.ListValidator)
	error = validator.Validate(*data, lv.CreateRules)
	if error != nil {
		response.ErrorWithMSG(fmt.Sprintf("保存失败了，%s", error.Error()))
		return
	}

	list, error := service.FindByID(data.Id)
	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	list, error = service.Update(list, data)
	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	response.SuccessWithData(*list)
}

func (ListController) Delete(request *gin.Context)  {
	response := response.Response{request}
	var error error

	form := service.NewModel()
	error = request.ShouldBindUri(form)
	if error != nil {
		response.ErrorWithMSG("删除失败")
		return
	}

	list, error := service.FindByID(form.Id)
	if error != nil {
		response.ErrorWithMSG("删除失败")
		return
	}

	error = service.Delete(list)

	if error != nil {
		response.ErrorWithMSG("删除失败")
		return
	}

	response.SuccessWithMSG("删除成功")
	return
}