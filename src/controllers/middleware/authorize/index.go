package authorize

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"todoList/src/models/user"
	"todoList/src/utils/response"
)


type Authorize struct {
}

func (Authorize) Auth(request *gin.Context)  {
	var response = response.Response{request}
	user := user.User()
	if user.Id == 0 {
		fmt.Println("check userInfo fail")
		response.ErrorWithMSG("请先登陆")
		request.Abort()
		return
	}

	request.Next()
}
