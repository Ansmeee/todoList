package list

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todoList/src/models/user"
	"todoList/src/services/list"
	"todoList/src/utils/response"
	lv "todoList/src/utils/validator/list"
)

type ListController struct {}
var service = new(list.ListService)

func (ListController) List(request *gin.Context)  {
	response := response.Response{request}
	var error error

	user := user.User()
	if len(user.Id) == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

	form := new(list.Params)
	error = request.ShouldBind(form)
	if error != nil {
		response.ErrorWithMSG("获取失败了，再试一次吧")
		return
	}

	form.Userid = user.Id
	total, data, error := service.List(form)
	if error != nil {
		response.ErrorWithMSG("获取失败了，再试一次吧")
		return
	}

	responseData := map[string]interface{}{
		"total": total,
		"list":  data,
	}

	response.SuccessWithData(responseData)
}

func (ListController) Create(request *gin.Context) {
	response := response.Response{request}
	var error error

	user := user.User()
	if len(user.Id) == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

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

	list.UserId = user.Id
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

	user := user.User()
	if len(user.Id) == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

	data := service.NewModel()
	error = request.ShouldBindUri(data)
	if error != nil {
		response.ErrorWithMSG("保存失败了，再试一次吧")
		return
	}

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

	user := user.User()
	if user.Id == "" {
		response.ErrorWithMSG("请先登陆")
		return
	}

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