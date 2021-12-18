package todo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todoList/src/services/list"
	todoService "todoList/src/services/todo"
	"todoList/src/utils/response"
	todoValidator "todoList/src/utils/validator/todo"
)

type TodoController struct{}

var thisService = &todoService.TodoService{}
var listService = &list.ListService{}

func (TodoController) List(request *gin.Context) {
	response := response.Response{request}
	var error error

	var form = new(todoService.QueryForm)
	error = request.ShouldBindQuery(form)
	if error != nil {
		response.ErrorWithMSG("请求失败，请重试")
		return
	}

	fmt.Println(form)

	if len(form.Rules) > 0 {
		var newRules [][]string
		for _, rule := range form.Rules {
			fmt.Println(rule)
			val := ""
			opt := "="
			if rule == "priority" {
				val = "高"
			}

			if rule == "status" {
				opt = "<>"
				val = "已完成"
			}

			newRules = append(newRules, []string{rule, opt, val})
		}

		form.Wheres = newRules
	}

	data, total, error := thisService.List(form)
	if error != nil {
		response.ErrorWithMSG("请求失败，请重试")
		return
	}

	responseData := map[string]interface{}{
		"list":  data,
		"total": total,
	}

	response.SuccessWithData(responseData)
	return
}

func (TodoController) Create(request *gin.Context) {
	var response = response.Response{request}

	var error error
	todo := thisService.NewModel()
	error = request.ShouldBind(todo)
	if error != nil {
		response.ErrorWithMSG("创建失败，请重试")
		return
	}

	validator := new(todoValidator.TodoValidator)
	error = validator.Validate(*todo, todoValidator.TodoCreateRules)
	if error != nil {
		response.ErrorWithMSG(fmt.Sprintf("创建失败: %s", error.Error()))
		return
	}

	if len(todo.ListId) > 0 {
		list, error := listService.FindByID(todo.ListId)
		if error != nil {
			response.ErrorWithMSG("创建失败，请重试")
			return
		}

		if len(list.Id) == 0 {
			response.ErrorWithMSG("创建失败，请重试")
			return
		}

		todo.Type = list.Type
	}

	data, error := thisService.Create(todo)
	if error != nil {
		response.ErrorWithMSG("创建失败，请重试")
		return
	}

	response.SuccessWithData(data)
}

func (TodoController) Detail(request *gin.Context) {
	var response = response.Response{request}
	var error error

	todo := thisService.NewModel()
	error = request.ShouldBindUri(todo)
	if error != nil || todo.Id == "" {
		response.ErrorWithMSG("获取失败，请重试")
		return
	}

	data, error := thisService.FindByID(todo.Id)
	if error != nil {
		response.ErrorWithMSG("获取失败，请重试")
		return
	}

	response.SuccessWithData(data)
}

func (TodoController) Update(request *gin.Context) {
	response := response.Response{request}
	var error error

	form := new(todoService.UpdateForm)
	error = request.ShouldBind(form)
	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	todo, error := thisService.FindByID(form.Id)
	if error != nil || todo.Id == "" {
		response.ErrorWithMSG("保存失败")
		return
	}

	error = thisService.Update(todo, form)
	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	response.SuccessWithData(*todo)
	return
}

func (TodoController) Delete(request *gin.Context) {
	response := response.Response{request}
	var error error

	form := thisService.NewModel()
	error = request.ShouldBindUri(form)
	if error != nil {
		response.ErrorWithMSG("删除失败")
		return
	}

	todo, error := thisService.FindByID(form.Id)
	if error != nil {
		response.ErrorWithMSG("删除失败")
		return
	}

	error = thisService.Delete(todo)
	if error != nil {
		response.ErrorWithMSG("删除失败")
		return
	}

	response.Success()
	return
}

type AttrForm struct {
	Id    string `form:"id"`
	Name  string `form:"name"`
	Value string `form:"value"`
}

func (TodoController) UpdateAttr(request *gin.Context) {
	response := response.Response{request}
	var error error

	attrForm := new(AttrForm)
	error = request.ShouldBind(attrForm)
	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	attrName := attrForm.Name
	attrValue := attrForm.Value

	todo, error := thisService.FindByID(attrForm.Id)
	if error != nil || todo.Id == "" {
		response.ErrorWithMSG("保存失败")
		return
	}

	if "" == attrName {
		response.ErrorWithMSG("保存失败")
		return
	}

	error = thisService.UpdateAttr(todo, attrName, attrValue)
	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	response.SuccessWithData(*todo)
	return
}
