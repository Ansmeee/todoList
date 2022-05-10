package user

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"path"
	"strings"
	cfg "todoList/config"
	userModel "todoList/src/models/user"
	"todoList/src/services/captcha"
	userService "todoList/src/services/user"
	"todoList/src/utils/response"
	userValidator "todoList/src/utils/validator/user"
)

type UserController struct {
}

var ctx = context.Background()
var service userService.UserService

func (UserController) CaptchaVerify(request *gin.Context) {
	id := request.Query("id")
	value := request.Query("value")

	captchaService := new(captcha.CaptchaService)
	fmt.Println(captchaService.Verify(id, value))
}

func (UserController) CaptchaId(request *gin.Context) {
	var response = response.Response{request}
	captchaService := new(captcha.CaptchaService)
	captchaID := captchaService.GenerateID()
	response.SuccessWithData(captchaID)
}

func (UserController) CaptchaImg(request *gin.Context) {
	captchaService := new(captcha.CaptchaService)
	source := request.Query("source")
	captchaService.GenerateImg(request.Writer, source)
}

type smsForm struct {
	Account string `form:"account"`
}

func (UserController) SendSMS(request *gin.Context)  {
	var response = response.Response{request}

	form := new(smsForm)
	err := request.ShouldBind(form)
	if err != nil {
		log.Println("SendSMS Error:", err)
		response.ErrorWithMSG("验证码发送失败")
		return
	}

	if form.Account == "" {
		log.Println("SendSMS Error: invalid account")
		response.ErrorWithMSG("验证码发送失败")
		return
	}

	existsCode := service.SMSCode(form.Account)
	if existsCode != "" {
		response.SuccessWithMSG("上条验证码依然有效")
		return
	}

	res := service.SendSMS(form.Account)
	if res {
		response.SuccessWithMSG("验证码已发送至您的手机")
		return
	}

	response.ErrorWithMSG("验证码发送失败")
	return
}

func (UserController) Info(request *gin.Context) {
	var response = response.Response{request}

	var user userModel.UserModel
	if err := request.ShouldBind(&user); err != nil {
		response.ErrorWithMSG("获取失败：参数错误")
		return
	}

	err, data := service.FindByID(user.Id)
	if err != nil {
		response.ErrorWithMSG(fmt.Sprintf("%s", err.Error()))
		return
	}

	if service.Verifing(data) {
		data.Verified = userModel.EMAIL_VERIFING
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

func (UserController) SignOut(request *gin.Context) {
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

func (UserController) SignIn(request *gin.Context) {
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
	//var response = response.Response{request}
	//
	//// 解析表单数据到 user model
	//var form userService.SignupForm
	//if err := request.ShouldBind(&form); err != nil {
	//	response.ErrorWithMSG("验证失败：参数错误")
	//	return
	//}
	//
	//// 参数验证
	//validator := new(userValidator.UserValidator)
	//if err := validator.Validate(form, userValidator.SignUpRules); err != nil {
	//	response.ErrorWithMSG(fmt.Sprintf("%s", err.Error()))
	//	return
	//}
	//
	//if form.Auth != form.PassWord {
	//	response.ErrorWithMSG(fmt.Sprintf("两次输入的密码不一致"))
	//	return
	//}
	//
	//// 注册
	//user, error := service.SignUp(&form)
	//if error != nil {
	//	response.ErrorWithMSG(fmt.Sprintf("%s", error.Error()))
	//	return
	//}
	//
	//token, err := service.GenerateToken(user)
	//if err != nil {
	//	response.SuccessWithDetail(302, "", nil)
	//	return
	//}
	//
	//res := service.LoginByToken(token, *user)
	//if res != true {
	//	response.SuccessWithDetail(302, "", nil)
	//	return
	//}
	//
	//var data = map[string]string{"token": token}
	//response.SuccessWithData(data)
}

func (UserController) ResetPass(request *gin.Context)  {
	response := response.Response{request}

	form := new(userService.ResetPassForm)
	error := request.ShouldBind(form)
	if error != nil {
		response.ErrorWithMSG("提交的信息有误，密码重置失败")
		return
	}

	error = service.ResetPass(form)
	if error != nil {
		response.ErrorWithMSG(error.Error())
		return
	}

	response.Success()
}

func (*UserController) EmailVerify(request *gin.Context)  {
	response := response.Response{request}

	form := map[string]string{"token": ""}
	if err := request.ShouldBind(&form); err != nil {
		response.ErrorWithMSG("验证失败，参数异常")
		return
	}

	token := form["token"]
	if token == "" {
		response.ErrorWithMSG("验证失败，token 不存在")
		return
	}

	res := service.Verified(token)
	if res == true {
		response.Success()
		return
	}

	response.ErrorWithMSG("验证失败")
	return
}

func (UserController) VerifyEmail(request *gin.Context)  {
	response := response.Response{request}

	user := userModel.User()
	if user.Id == "" {
		response.ErrorWithMSG("请先登陆")
		return
	}

	if userModel.EMAIL_VERIFIED == user.Verified {
		response.ErrorWithMSG("邮箱已完成验证")
		return
	}

	if service.Verifing(user) == true {
		response.ErrorWithMSG("邮箱正在验证中")
		return
	}

	error := service.VerifyEmail(user)
	if error != nil {
		response.ErrorWithMSG("验证失败")
		return
	}

	response.Success()
	return
}

func (UserController) UpdateAttr(request *gin.Context) {
	response := response.Response{request}

	user := userModel.User()
	if user.Id == "" {
		response.ErrorWithMSG("请先登陆")
		return
	}

	form := new(userService.AttrForm)
	error := request.ShouldBind(form)
	if error != nil {
		response.ErrorWithMSG("更新失败")
		return
	}

	if user.Id != form.Id {
		response.ErrorWithMSG("更新失败")
		return
	}

	attrKey := form.Key
	attrVal := form.Value

	if "" == attrKey {
		response.ErrorWithMSG("更新失败")
		return
	}

	if attrKey == "phone"{
		if form.Code == "" {
			response.ErrorWithMSG("请输入短信验证码")
			return
		}

		code := service.SMSCode(attrVal)
		if code == "" || code != form.Code {
			response.ErrorWithMSG("请输入正确的短信验证码")
			return
		}
	}

	error = service.UpdateAttr(user, attrKey, attrVal)
	if error != nil {
		response.ErrorWithMSG("更新失败")
		return
	}

	response.SuccessWithData(*user)

}

func (UserController) Update(request *gin.Context) {
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

	if err := service.Update(user, &updateUser); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("操作失败：%s", err.Error()))
	}

	response.Success()
}

func (UserController) Delete(request *gin.Context) {
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

	if err := service.Delete(user); err != nil {
		response.ErrorWithMSG(fmt.Sprintf("删除失败：%s", err.Error()))
		return
	}
	response.Success()
}

func (UserController) ShowIcon(request *gin.Context)  {
	response := response.Response{request}

	icon := request.Param("icon")

	iconPath, error := service.GenerateLocalIconPath()
	if error != nil {
		response.ErrorWithMSG("获取失败")
	}

	file := fmt.Sprintf("%s/%s", iconPath, icon)
	_, error = os.Stat(file)
	if error != nil {
		response.ErrorWithMSG("头像加载失败")
		return
	}

	response.SuccessWithFile(file)
}

func (UserController) Icon(request *gin.Context) {
	response := response.Response{request}

	user := userModel.User()
	if user.Id == "" {
		response.ErrorWithMSG("请先登陆")
		return
	}

	form := new(userService.AttrForm)
	error := request.ShouldBind(form)

	if error != nil {
		response.ErrorWithMSG("上传失败")
		return
	}

	if error != nil || user.Id != form.Id {
		response.ErrorWithMSG("上传失败")
		return
	}

	file, fileHeader, error := request.Request.FormFile("icon")
	if error != nil {
		response.ErrorWithMSG("上传失败")
		return
	}

	if fileHeader.Size > 1024 * 1024 * 2 {
		response.ErrorWithMSG("图片太大了")
		return
	}

	enableImgTypes := map[string]bool{".png": true, ".jpg": true, ".jpeg": true}
	ext := path.Ext(fileHeader.Filename)
	if _, ok := enableImgTypes[ext]; !ok {
		response.ErrorWithMSG("只能上传 .png .jpg .jpeg 类型的图片")
		return
	}

	conf, error := cfg.Config()
	if error != nil {
		response.ErrorWithMSG("上传失败")
		return
	}

	saveHandler := conf.Section("environment").Key("icon_save_handler").String()

	var url = ""
	if "QN" == saveHandler {
		url, error = service.SaveIcon2QN(file, fileHeader)
	} else {
		iconPath, error := service.GenerateLocalIconPath()
		if error != nil {
			response.ErrorWithMSG("上传失败")
		}

		name := fmt.Sprintf("%x", md5.Sum([]byte(user.Id)))
		fileName := fmt.Sprintf("%s%s", name, ext)
		savePath := fmt.Sprintf("/%s/%s", strings.Trim(iconPath, "/"), fileName)

		error = request.SaveUploadedFile(fileHeader, savePath)

		host := conf.Section("environment").Key("app_host").String()
		url = fmt.Sprintf("%s/rest/user/icon/%s", host, fileName)
	}

	if error != nil {
		fmt.Println(error.Error())
		response.ErrorWithMSG("上传失败")
	}

	error = service.UpdateAttr(user, "icon", url)
	if error != nil {
		fmt.Println(error.Error())
		response.ErrorWithMSG("上传失败")
	}

	response.SuccessWithData(url)
}
