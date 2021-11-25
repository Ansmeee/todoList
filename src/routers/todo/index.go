package todo

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/todo"
)

type TodoRouter struct {
}

func (*TodoRouter) InitRouter(group *gin.RouterGroup) {
	//middleware := new(authorize.Authorize)
	router := group.Group("todo")
	controller := new(todo.TodoController)
	{
		router.GET("list", controller.List)
		router.POST("", controller.Create)
		router.GET("", controller.Detail)
		router.PUT("", controller.Update)
	}
}
