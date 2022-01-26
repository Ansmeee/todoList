package list

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/list"
	"todoList/src/controllers/middleware/authorize"
)

type ListRouter struct{}

func (ListRouter) InitRouter(router gin.IRoutes) {
	auth := new(authorize.Authorize)
	router = router.Use(auth.Auth)

	controller := new(list.ListController)
	router.GET("list", controller.List)
	router.POST("list", controller.Create)
	router.PUT("list/:id", controller.Update)
	router.DELETE("list/:id", controller.Delete)

}
