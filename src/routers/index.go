package routers

import (
	"github.com/gin-gonic/gin"
	"todoList/src/routers/feedback"
	"todoList/src/routers/list"
	"todoList/src/routers/msg"
	"todoList/src/routers/static"
	"todoList/src/routers/todo"
	"todoList/src/routers/user"
)

// 初始化，注册路由
type RouterGroup struct {
	StaticRouter   static.StaticRouter
	UserRouter     user.UserRouter
	TodoRouter     todo.TodoRouter
	ListRouter     list.ListRouter
	MsgRouter      msg.MsgRouter
	FeedbackRouter feedback.FeedbackRouter
}

func (group *RouterGroup) InitRouter(routers gin.IRoutes) {
	group.StaticRouter.InitRouter(routers)
	group.UserRouter.InitRouter(routers)
	group.TodoRouter.InitRouter(routers)
	group.ListRouter.InitRouter(routers)
	group.MsgRouter.InitRouter(routers)
	group.FeedbackRouter.InitRouter(routers)
}
