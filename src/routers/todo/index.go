package todo

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/middleware/authorize"
	"todoList/src/controllers/todo"
)

type TodoRouter struct {
}

func (*TodoRouter) InitRouter(router gin.IRoutes) {
	controller := new(todo.TodoController)
	auth := new(authorize.Authorize)
	arouter := router.Use(auth.Auth)
	arouter.GET("todo", controller.List)
	arouter.POST("todo", controller.Create)
	arouter.PUT("todo", controller.Update)
	arouter.GET("todo/:id", controller.Detail)
	arouter.DELETE("todo/:id", controller.Delete)
	arouter.PUT("todo/attr", controller.UpdateAttr)
	arouter.POST("todo/upload", controller.Upload)
}
