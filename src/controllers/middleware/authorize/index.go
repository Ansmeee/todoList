package authorize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todoList/src/models/user"
	userService "todoList/src/services/user"
	"todoList/src/utils/response"
)


type Authorize struct {
}


func (Authorize) Auth(request *gin.Context)  {
	var response = response.Response{request}

	token := request.GetHeader("token")

	userService := new(userService.UserService)
	userInfo, err := userService.GetUserInfoByToken(token)
	if err != nil {
		fmt.Println("GetUserInfoByToken 失败")
		response.ErrorWithMSG("获取失败：非法的请求")
		request.Abort()
		return
	}

	if userInfo.Id == "" {
		fmt.Println("check userInfo 失败")
		response.ErrorWithMSG("获取失败：用户信息异常")
		request.Abort()
		return
	}

	if !userInfo.OnJob() {
		fmt.Println("check user on job 失败")
		response.ErrorWithMSG("获取失败：用户信息异常")
		request.Abort()
		return
	}

	authModel := new(user.AuthModel)
	authModel.Login(&userInfo)

	request.Next()
}
