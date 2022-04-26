package user

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/middleware/authorize"
	"todoList/src/controllers/user"
)

type UserRouter struct {
}

func (*UserRouter) InitRouter(router gin.IRoutes) {
	auth := new(authorize.Authorize)
	authRouter := router.Use(auth.Auth)
	controller := new(user.UserController)
	authRouter.POST("user/icon", controller.Icon)
	authRouter.GET("user", controller.Info)
	authRouter.DELETE("user", controller.Delete)
	authRouter.PUT("user", controller.Update)
	authRouter.PUT("user/attr", controller.UpdateAttr)
	authRouter.PUT("user/pass", controller.ResetPass)
	authRouter.POST("user/verify/email", controller.VerifyEmail)
}
