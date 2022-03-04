package msg

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/middleware/authorize"
	"todoList/src/controllers/msgController"
)

type MsgRouter struct {
}

func (*MsgRouter) InitRouter(router gin.IRoutes) {
	auth := new(authorize.Authorize)
	router = router.Use(auth.Auth)

	controller := new(msgController.MsgController)
	router.GET("msg", controller.List)
	router.PUT("msg/:id", controller.UpdateAttr)
}
