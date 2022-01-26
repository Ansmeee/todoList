package list

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/list"
)

type ListRouter struct{}

func (ListRouter) InitRouter(router gin.IRoutes) {
	controller := new(list.ListController)
	{
		router.GET("list", controller.List)
		router.POST("list", controller.Create)
		router.PUT("list/:id", controller.Update)
		router.DELETE("list/:id", controller.Delete)
	}
}
