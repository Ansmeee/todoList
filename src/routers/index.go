package routers

import (
	"github.com/gin-gonic/gin"
	"todoList/src/routers/todo"
	"todoList/src/routers/user"
)

// 初始化，注册路由
type RouterGroup struct {
	UserRouter user.UserRouter
	TodoRouter todo.TodoRouter
}

func (group *RouterGroup) InitRouter(routerGroup *gin.RouterGroup) {
	group.UserRouter.InitUserRouter(routerGroup)
	group.TodoRouter.InitTodoRouter(routerGroup)
}