package user

import (
	"todoList/src/controllers/middleware/authorize"
	"todoList/src/controllers/user"
	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter (group *gin.RouterGroup)  {
	middleware := new(authorize.Authorize)
	router := group.Group("user")
	controller := new (user.UserController)
	{
		router.POST("signin", controller.SignIn)
		router.POST("signup", controller.SignUp)
		router.POST("delete", controller.Delete).Use(middleware.Auth)
		router.POST("update", controller.Update).Use(middleware.Auth)
		router.GET("info/:id", controller.Info).Use(middleware.Auth)
		router.GET("list", controller.List).Use(middleware.Auth)
	}
}

