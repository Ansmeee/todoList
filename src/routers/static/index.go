package static

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/todo"
	"todoList/src/controllers/user"
)

type StaticRouter struct {
}

func (*StaticRouter) InitRouter(router gin.IRoutes) {
	userCTRL := new(user.UserController)
	router.GET("user/icon/:icon", userCTRL.ShowIcon)
	router.GET("user/captchaimg", userCTRL.CaptchaImg)
	router.GET("user/captchaid", userCTRL.CaptchaId)
	router.GET("user/captchaverify", userCTRL.CaptchaVerify)
	router.POST("user/signin", userCTRL.SignIn)
	router.POST("user/signout", userCTRL.SignOut)
	router.POST("user/signup", userCTRL.SignUp)
	router.POST("user/email/verify", userCTRL.EmailVerify)
	router.POST("user/sms", userCTRL.SendSMS)

	todoCTRL := new(todo.TodoController)
	router.GET("todo/:id/img/:img", todoCTRL.Img)
}