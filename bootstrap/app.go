package bootstrap

import (
	"github.com/gin-gonic/gin"
	"todoList/src/controllers/middleware"
	"todoList/src/routers"
)

func InitEngine() (engine *gin.Engine) {
	engine = gin.Default()
	routerGroup := new (routers.RouterGroup)
	var routers = engine.Group("rest").Use(middleware.Auth)
	routerGroup.InitRouter(routers)

	return
}