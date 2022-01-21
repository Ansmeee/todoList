package user

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/user"
)

type UserRouter struct {
}

func (*UserRouter) InitRouter(group *gin.RouterGroup) {
	router := group.Group("user")
	controller := new(user.UserController)
	{
		router.POST("signin", controller.SignIn)
		router.POST("signout", controller.SignOut)
		router.POST("signup", controller.SignUp)
		router.POST("icon", controller.Icon)
		router.GET("", controller.Info)
		router.GET("list", controller.List)
		router.DELETE("", controller.Delete)
		router.PUT("", controller.Update)

	}
}
