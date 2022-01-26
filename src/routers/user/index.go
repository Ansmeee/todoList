package user

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/user"
)

type UserRouter struct {
}

func (*UserRouter) InitRouter(router gin.IRoutes) {
	controller := new(user.UserController)
	router.POST("user/signin", controller.SignIn)
	router.POST("user/signout", controller.SignOut)
	router.POST("user/signup", controller.SignUp)
	router.POST("user/icon", controller.Icon)
	router.GET("user", controller.Info)
	router.GET("user/list", controller.List)
	router.DELETE("user", controller.Delete)
	router.PUT("user", controller.Update)
	router.PUT("user/attr", controller.UpdateAttr)
}
