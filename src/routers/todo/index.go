package todo

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/todo"
)

type TodoRouter struct {
}

func (*TodoRouter) InitRouter(router gin.IRoutes) {
	controller := new(todo.TodoController)
	router.GET("todo", controller.List)
	router.POST("todo", controller.Create)
	router.PUT("todo", controller.Update)
	router.GET("todo/:id", controller.Detail)
	router.DELETE("todo/:id", controller.Delete)
	router.PUT("todo/attr", controller.UpdateAttr)
}
