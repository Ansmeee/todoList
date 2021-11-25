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
		router.POST("", controller.Create)
	}
}