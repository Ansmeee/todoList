package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"path/filepath"
	userModel "todoList/src/models/user"
	userService "todoList/src/services/user"
	"todoList/src/utils/response"
	userValidator "todoList/src/utils/validator/user"
)

type UserController struct {
}

var service userService.UserService

func (UserController) Info(request *gin.Context) {
	var response = response.Response{request}

	var user userModel.UserModel
	if err := request.ShouldBind(&user); err != nil {
		response.ErrorWithMSG("获取失败：参数错误")
		return
	}

	err, data := service.FindByID(user.Id)
	if err != nil {
		response.ErrorWithMSG(fmt.Sprintf("获取失败：%s", err.Error()))
		return
	}

	response.SuccessWithData(data)
	return
}

func (UserController) List(request *gin.Context) {
	var response = response.Response{request}

	var form = new(userService.QueryParams)
	if err := request.ShouldBindQuery(&form); err != nil {
		response.ErrorWithMSG("获取失败：参数异常")
		return
	}

	err, data, total := service.List(form)
	if err != nil {
		response.ErrorWithMSG(fmt.Sprintf("获取失败：%s", err.Error()))
		return
	}

	var responseData = map[string]interface{}{
		"list":  data,
		"total": total,
	}

	response.SuccessWithData(responseData)
	return
}

func (UserController) SignOut(request *gin.Context)  {
	var response = response.Response{request}

	token := request.GetHeader("Authorization")
	if len(token) == 0 {
		response.Success()
		return
	}

	error := service.SignOut(token)
	if error != nil {
		response.ErrorWithMSG(fmt.Sprintf("登出失败，%s", error.Error()))
		return
	}

	response.Success()
}

func (UserController) SignIn(request *gin.Context)  {
	var response = response.Response{request}

	form := new(userService.SigninForm)
	if err := request.ShouldBind(form); err != nil {
		response.ErrorWithMSG("登录失败")
		return
	}

	validator := new(userValidator.UserValidator)
	if err := validator.Validate(*form, userValidator.SignInRules); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("登录失败，%s", err.Error()))
		return
	}

	token, err := service.SignIn(form)
	if err != nil {
		response.ErrorWithMSG(fmt.Sprintf("登录失败，%s", err.Error()))
		return
	}

	var data = map[string]string{"token": token}
	response.SuccessWithData(data)
}

func (UserController) SignUp(request *gin.Context) {
	var response = response.Response{request}

	// 解析表单数据到 user model
	var form userService.SignupForm
	if err := request.ShouldBind(&form); err != nil {
		response.ErrorWithMSG("验证失败：参数错误")
		return
	}

	// 参数验证
	validator := new(userValidator.UserValidator)
	if err := validator.Validate(form, userValidator.SignUpRules); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("验证失败：%s", err.Error()))
		return
	}

	err, existUser := service.FindeByEmail(form.Email)
	fmt.Println(existUser)
	if err != nil || existUser.Id != 0 {
		response.ErrorWithMSG(fmt.Sprintf("验证失败：无法注册该账号"))
		return
	}

	fmt.Println(form.Auth, form.PassWord)
	if form.Auth != form.PassWord {
		response.ErrorWithMSG(fmt.Sprintf("验证失败：两次输入的密码不一致"))
		return
	}

	// 注册
	if err := service.SignUp(&form); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("创建失败：%s", err.Error()))
		return
	}

	response.Success()
	return
}

func (UserController) Update(request *gin.Context)  {
	response := response.Response{request}

	var updateUser userModel.UserModel
	if err := request.ShouldBind(&updateUser); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("操作失败：参数异常"))
		return
	}

	err, user := service.FindByID(updateUser.Id)
	if err != nil {
		response.ErrorWithMSG("操作失败：该用户不存在")
		return
	}

	if err := service.Update(&user, &updateUser); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("操作失败：%s", err.Error()))
	}

	response.Success()
}

func (UserController) Delete(request *gin.Context)  {
	response := response.Response{request}

	userForm := new(userModel.UserModel)
	if err := request.ShouldBind(&userForm); err != nil {
		response.ErrorWithMSG("删除失败：参数异常")
		return
	}

	err, user := service.FindByID(userForm.Id)
	if err != nil {
		response.ErrorWithMSG("删除失败：用户不存在")
		return
	}

	if err := service.Delete(&user); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("删除失败：%s", err.Error()))
		return
	}
	response.Success()
}

func (UserController) Icon(request *gin.Context)  {
	response := response.Response{request}

	file, error := request.FormFile("icon")
	if error != nil {
		response.ErrorWithMSG("上传失败")
		return
	}

	basePath := "./"
	filename := basePath + filepath.Base(file.Filename)
	if error = request.SaveUploadedFile(file, filename); error != nil {
		fmt.Println(error.Error())
		response.ErrorWithMSG("上传失败")
		return
	}

	response.Success()
}