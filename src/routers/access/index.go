package access

import (
	"todoList/src/controllers/access"
	"todoList/src/controllers/middleware/authorize"
	"github.com/gin-gonic/gin"
)

type AccessRouter struct {
}

func (*AccessRouter) InitUserRouter(group *gin.RouterGroup)  {
	middleware := new(authorize.Authorize)
	router := group.Group("access").Use(middleware.Auth)
	controller := new(access.AccessController)
	{
		router.POST("createRole", controller.CreateRole)
		router.GET("roleList", controller.RoleList)
	}

}