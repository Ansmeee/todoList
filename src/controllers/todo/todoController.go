package todo

import (
	"fmt"
	"github.com/gin-gonic/gin"
	todoService "todoList/src/services/todo"
	"todoList/src/utils/response"
	todoValidator "todoList/src/utils/validator/todo"
)

type TodoController struct {
}

var thisService = &todoService.TodoService{}
func init()  {
	thisService = new(todoService.TodoService)
}

func (TodoController) List(request *gin.Context)  {

}

func (TodoController) Create(request *gin.Context)  {
	var response = response.Response{request}

	var error error
	todo := thisService.NewModel()
	error = request.ShouldBind(todo)
	if error != nil {
		response.ErrorWithMSG("创建失败，请重试")
		return
	}

	validator := new(todoValidator.TodoValidator)
	error = validator.Validate(todo, todoValidator.TodoCreateRules)
	if error != nil {
		response.ErrorWithMSG(fmt.Sprintf("创建失败: %s", error.Error()))
		return
	}


	data, error := thisService.Create(todo)
	if error != nil {
		response.ErrorWithMSG("创建失败，请重试")
		return
	}

	response.SuccessWithData(data)
}

func (TodoController) Detail(request *gin.Context)  {
	var response = response.Response{request}
	var error error

	todo := thisService.NewModel()
	error = request.ShouldBind(todo)
	if error != nil {
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

func (TodoController) Update(request *gin.Context)  {

}