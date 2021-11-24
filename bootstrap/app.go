package bootstrap

import (
	"todoList/src/routers"
	"github.com/gin-gonic/gin"
)

func InitEngine() (engine *gin.Engine) {
	//redis.Client = redis.Connect()
	//database.DB = database.Connect("default")

	engine = gin.Default()
	routerGroup := new (routers.RouterGroup)
	var router = engine.Group("")
	{routerGroup.InitRouter(router)}

	return
}