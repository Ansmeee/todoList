package routers

import (
	"todoList/src/routers/access"
	"todoList/src/routers/user"
	"github.com/gin-gonic/gin"
)

// 初始化，注册路由
type RouterGroup struct {
	UserRouter user.UserRouter
	AccessRouter access.AccessRouter
}

func (group *RouterGroup) InitRouter(routerGroup *gin.RouterGroup) {
	group.UserRouter.InitUserRouter(routerGroup)
	group.AccessRouter.InitUserRouter(routerGroup)
}