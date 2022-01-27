package todo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
	"todoList/src/models/todo"
	"todoList/src/models/user"
	"todoList/src/services/list"
	todoService "todoList/src/services/todo"
	"todoList/src/utils/response"
	todoValidator "todoList/src/utils/validator/todo"
)

type TodoController struct{}

var thisService = &todoService.TodoService{}
var listService = &list.ListService{}

func setLatestForm(request *gin.Context, form *todoService.QueryForm) {
	form.PageSize = 20

	sortBy := "updated_at"
	if len(form.SortBy) > 0 {
		sortBy = form.SortBy
	}
	form.SortBy = sortBy

	sortOrder := "desc"
	if len(form.SortOrder) > 0 {
		sortOrder = form.SortOrder
	}
	form.SortOrder = sortOrder

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

	newRules = append(newRules, []string{"status", "<=", status}, []string{"user_id", "=", strconv.Itoa(user.User().Id)})
	form.Wheres = newRules
}

func setDoneForm(request *gin.Context, form *todoService.QueryForm) {
	sortBy := "finished_at"
	if len(form.SortBy) > 0 {
		sortBy = form.SortBy
	}
	form.SortBy = sortBy

	sortOrder := "desc"
	if len(form.SortOrder) > 0 {
		sortOrder = form.SortOrder
	}
	form.SortOrder = sortOrder

	var newRules [][]string
	form.Rules = request.QueryArray("rules[]")
	if len(form.Rules) > 0 {
		for _, rule := range form.Rules {
			val := ""
			opt := "="
			if rule == "priority" {
				val = "3"
			}

			newRules = append(newRules, []string{rule, opt, val})
		}
	}

	newRules = append(newRules, []string{"status", "=", "2"}, []string{"user_id", "=", strconv.Itoa(user.User().Id)})
	form.Wheres = newRules
}

func setDefaultForm(request *gin.Context, form *todoService.QueryForm) {
	form.ListId = form.From

	sortBy := "created_at"
	if len(form.SortBy) > 0 {
		sortBy = form.SortBy
	}
	form.SortBy = sortBy

	sortOrder := "desc"
	if len(form.SortOrder) > 0 {
		sortOrder = form.SortOrder
	}
	form.SortOrder = sortOrder

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

	newRules = append(newRules, []string{"status", "<=", status}, []string{"user_id", "=", strconv.Itoa(user.User().Id)})
	form.Wheres = newRules
}

func (TodoController) List(request *gin.Context) {
	response := response.Response{request}
	var error error

	user := user.User()
	if user.Id == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

	var form = new(todoService.QueryForm)
	error = request.ShouldBindQuery(form)
	if error != nil {
		response.ErrorWithMSG("请求失败，请重试")
		return
	}

	switch form.From {
	case "done":
		setDoneForm(request, form)
	case "latest":
		setLatestForm(request, form)
	default:
		setDefaultForm(request, form)
	}

	data, total, error := thisService.List(form)
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

	user := user.User()
	if user.Id == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

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

	todo.UserId = user.Id

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

	user := user.User()
	if user.Id == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

	todo := thisService.NewModel()
	error = request.ShouldBindUri(todo)
	if error != nil || todo.Id == 0 {
		response.ErrorWithMSG("获取失败，请重试")
		return
	}

	data, error := thisService.FindByID(todo.Id)
	if error != nil || data.UserId != user.Id {
		response.ErrorWithMSG("获取失败，请重试")
		return
	}

	response.SuccessWithData(data)
}

func (TodoController) Update(request *gin.Context) {
	response := response.Response{request}
	var error error

	user := user.User()
	if user.Id == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

	form := thisService.NewModel()
	error = request.ShouldBind(form)

	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	todo, error := thisService.FindByID(form.Id)
	if error != nil || todo.Id == 0 || todo.UserId != user.Id {
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
	user := user.User()
	if user.Id == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

	form := thisService.NewModel()
	error := request.ShouldBindUri(form)
	if error != nil {
		response.ErrorWithMSG("删除失败")
		return
	}

	todo, error := thisService.FindByID(form.Id)
	if error != nil || todo.UserId != user.Id{
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

	user := user.User()
	if user.Id == 0 {
		response.ErrorWithMSG("请先登陆")
		return
	}

	attrForm := new(AttrForm)
	error := request.ShouldBind(attrForm)
	if error != nil {
		response.ErrorWithMSG("保存失败")
		return
	}

	attrName := attrForm.Name
	attrValue := attrForm.Value

	todo, error := thisService.FindByID(attrForm.Id)
	if error != nil || todo.Id == 0 || todo.UserId != user.Id{
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
