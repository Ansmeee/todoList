package user

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/middleware/authorize"
	"todoList/src/controllers/user"
)

type UserRouter struct {
}

func (*UserRouter) InitRouter(group *gin.RouterGroup) {
	middleware := new(authorize.Authorize)
	router := group.Group("user")
	controller := new(user.UserController)
	{
		router.POST("signin", controller.SignIn)
		router.POST("signout", controller.SignOut)
		router.POST("signup", controller.SignUp)
		router.GET("list", controller.List).Use(middleware.Auth)
		router.DELETE("", controller.Delete).Use(middleware.Auth)
		router.PUT("", controller.Update).Use(middleware.Auth)
		router.GET("", controller.Info).Use(middleware.Auth)
	}
}
