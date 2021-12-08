package list

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/list"
)

type ListRouter struct {}

func (ListRouter) InitRouter(group *gin.RouterGroup)  {
	router := group.Group("list")
	controller := new(list.ListController)
	{
		router.GET("", controller.List)
		router.POST("", controller.Create)
		router.PUT("/:id", controller.Update)
		router.DELETE("/:id", controller.Delete)
	}
}