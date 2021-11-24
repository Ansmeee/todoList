package bootstrap

import (
	"github.com/gin-gonic/gin"
	"todoList/src/routers"
)

func InitEngine() (engine *gin.Engine) {
	engine = gin.Default()
	routerGroup := new (routers.RouterGroup)
	var router = engine.Group("")
	{routerGroup.InitRouter(router)}

	return
}