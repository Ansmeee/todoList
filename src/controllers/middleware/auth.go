package middleware

import (
	"github.com/gin-gonic/gin"
	"todoList/src/models/user"
	userService "todoList/src/services/user"
)

func Auth(request *gin.Context)  {
	token := request.GetHeader("Authorization")

	userInfo := new(user.UserModel)
	authModel := new(user.AuthModel)
	if len(token) == 0 {
		authModel.SetUser(userInfo)
		request.Next()
		return
	}

	userService := new(userService.UserService)
	userInfo, _ = userService.GetUserInfoByToken(token)
	authModel.SetUser(userInfo)

	request.Next()
}
