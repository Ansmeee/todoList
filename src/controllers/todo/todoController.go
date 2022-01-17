package todo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
	"todoList/src/models/todo"
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

	var form = todoService.QueryForm{From: "", SortBy: "created_at", SortOrder: "desc"}
	error = request.ShouldBindQuery(&form)
	if error != nil {
		response.ErrorWithMSG("请求失败，请重试")
		return
	}

	if form.From == "latest" {
		form.PageSize = 20
	} else {
		form.ListId = form.From
	}

	var status = "1"
	var newRules [][]string
	form.Rules = request.QueryArray("rules[]")
	if len(form.Rules) > 0 {
		for _, rule := range form.Rules {
			val := ""
			opt := "="
			if rule == "priority" {
				val = "3"
			}

			if rule == "status" {
				status = "2"
				continue
			}

			newRules = append(newRules, []string{rule, opt, val})
		}
	}

	newRules = append(newRules, []string{"status", "<=", status})
	form.Wheres = newRules

	data, total, error := thisService.List(&form)
	if error != nil {
		response.ErrorWithMSG("请求失败，请重试")
		return
	}

	if len(data) == 0 {
		data = []todo.TodoModel{}
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

	if todo.ListId != 0 {
		list, error := listService.FindByID(todo.ListId)
		if error != nil {
			response.ErrorWithMSG("创建失败，请重试")
			return
		}

		if list.Id == 0 {
			response.ErrorWithMSG("创建失败，请重试")
			return
		}

		todo.ListId = list.Id
		todo.Type = list.Type
	}

	if len(todo.Title) == 0 {
		todo.Title = "未命名"
	}

	if len(todo.Deadline) == 0 {
		todo.Deadline = time.Now().Format("2006-01-02")
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
	if error != nil || todo.Id == 0 {
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

	form :=  thisService.NewModel()
	error = request.ShouldBind(form)

	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	todo, error := thisService.FindByID(form.Id)
	if error != nil || todo.Id == 0 {
		response.ErrorWithMSG("保存失败")
		return
	}

	if len(form.Title) == 0 {
		form.Title = "未命名"
	}

	if form.ListId != 0 {
		list, error := listService.FindByID(form.ListId)
		if error != nil {
			response.ErrorWithMSG("更新失败")
			return
		}

		if list.Id == 0 {
			response.ErrorWithMSG("更新失败")
			return
		}

		form.ListId = list.Id
		form.Type = list.Title
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
	Id    int    `form:"id"`
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
	if error != nil || todo.Id == 0 {
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
