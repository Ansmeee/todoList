package middleware

import (
	"github.com/gin-gonic/gin"
	"todoList/src/models/user"
	userService "todoList/src/services/user"
)

func Auth(request *gin.Context)  {
	token := request.GetHeader("Authorization")
	if len(token) == 0 {
		request.Next()
		return
	}

	userService := new(userService.UserService)
	userInfo, err := userService.GetUserInfoByToken(token)
	if err != nil {
		request.Next()
		return
	}

	if userInfo.Id == 0 {
		request.Next()
		return
	}

	authModel := new(user.AuthModel)
	authModel.SetUser(userInfo)

	request.Next()
}
