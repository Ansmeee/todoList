package routers

import (
	"github.com/gin-gonic/gin"
	"todoList/src/routers/list"
	"todoList/src/routers/msg"
	"todoList/src/routers/todo"
	"todoList/src/routers/user"
)

// 初始化，注册路由
type RouterGroup struct {
	UserRouter user.UserRouter
	TodoRouter todo.TodoRouter
	ListRouter list.ListRouter
	MsgRouter  msg.MsgRouter
}

func (group *RouterGroup) InitRouter(routers gin.IRoutes) {
	group.UserRouter.InitRouter(routers)
	group.TodoRouter.InitRouter(routers)
	group.ListRouter.InitRouter(routers)
	group.MsgRouter.InitRouter(routers)
}
