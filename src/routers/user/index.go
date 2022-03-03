package user

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/middleware/authorize"
	"todoList/src/controllers/user"
)

type UserRouter struct {
}

func (*UserRouter) InitRouter(router gin.IRoutes) {
	controller := new(user.UserController)

	router.GET("user/icon/:icon", controller.ShowIcon)
	router.GET("user/captchaimg", controller.CaptchaImg)
	router.GET("user/captchaid", controller.CaptchaId)
	router.GET("user/captchaverify", controller.CaptchaVerify)
	router.POST("user/signin", controller.SignIn)
	router.POST("user/signout", controller.SignOut)
	router.POST("user/signup", controller.SignUp)

	auth := new(authorize.Authorize)
	authRouter := router.Use(auth.Auth)
	authRouter.POST("user/icon", controller.Icon)
	authRouter.GET("user", controller.Info)
	authRouter.DELETE("user", controller.Delete)
	authRouter.PUT("user", controller.Update)
	authRouter.PUT("user/attr", controller.UpdateAttr)
}
