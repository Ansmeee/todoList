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
		router.GET("", controller.List)
		router.POST("", controller.Create)
		router.PUT("", controller.Update)
		router.GET("/:id", controller.Detail)
		router.DELETE("/:id", controller.Delete)
		router.PUT("item", controller.Item)
	}
}
