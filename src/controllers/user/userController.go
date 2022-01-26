package user

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	userModel "todoList/src/models/user"
	userService "todoList/src/services/user"
	"todoList/src/utils/response"
	userValidator "todoList/src/utils/validator/user"
)

type UserController struct {
}

var ctx = context.Background()
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
		response.ErrorWithMSG(fmt.Sprintf("%s", err.Error()))
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
		response.ErrorWithMSG(fmt.Sprintf("%s", err.Error()))
		return
	}

	if form.Way == "email" {
		err, existUser := service.FindeByEmail(form.Account)
		if err != nil || existUser.Id != 0 {
			response.ErrorWithMSG(fmt.Sprintf("该邮箱已被占用"))
			return
		}
	}

	if form.Way == "phone" {

	}

	if form.Auth != form.PassWord {
		response.ErrorWithMSG(fmt.Sprintf("两次输入的密码不一致"))
		return
	}

	// 注册
	user, error := service.SignUp(&form)
	if  error != nil {
		response.ErrorWithMSG(fmt.Sprintf("%s", error.Error()))
		return
	}

	token, err := service.GenerateToken(user)
	if err != nil {
		response.SuccessWithDetail(302, "", nil)
		return
	}

	var data = map[string]string{"token": token}
	response.SuccessWithData(data)
}

func (UserController) UpdateAttr(request *gin.Context) {
	response := response.Response{request}

	form := new(userService.AttrForm)
	error := request.ShouldBind(form)
	if error != nil {
		response.ErrorWithMSG("更新失败")
		return
	}

	error, user := service.FindByID(form.Id)
	if error != nil {
		response.ErrorWithMSG("更新失败")
		return
	}

	if user.Id == 0 {
		response.ErrorWithMSG("更新失败，用户信息异常")
		return
	}

	attrKey := form.Key
	attrVal := form.Value

	if "" == attrKey {
		response.ErrorWithMSG("更新失败")
		return
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

func (UserController) Icon(request *gin.Context) {
	response := response.Response{request}

	form := new(userService.AttrForm)
	error := request.ShouldBind(form)

	if error != nil {
		response.ErrorWithMSG("上传失败")
		return
	}

	error, user := service.FindByID(form.Id)
	if error != nil || user.Id == 0 {
		response.ErrorWithMSG("上传失败")
		return
	}

	file, fileHeader, error := request.Request.FormFile("icon")
	if error != nil {
		response.ErrorWithMSG("上传失败")
		return
	}

	putPolicy := storage.PutPolicy{Scope: "ansmetodolist"}
	mac := qbox.NewMac("grcCCRTRuJwq0OKb4VUbxZm5L2_FQlJyUex0mN85", "EYd181-l6Rc5yvGNtszmHvurp9qiaYsfGgVktF5f")
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          &storage.ZoneHuabei,
		UseCdnDomains: false,
	}

	imgHost := "http://r646b3gyv.hb-bkt.clouddn.com"
	fileSize := fileHeader.Size
	putExtra := storage.PutExtra{}
	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	err := formUploader.PutWithoutKey(ctx, &ret, upToken, file, fileSize, &putExtra)
	if err != nil {
		fmt.Println(err.Error())
		response.ErrorWithMSG("上传失败")
		return
	}

	url := storage.MakePublicURL(imgHost, ret.Key)
	error = service.UpdateAttr(user, "icon", url)
	if error !=nil {
		fmt.Println(error.Error())
		response.ErrorWithMSG("上传失败")
		return
	}

	response.SuccessWithData(url)
}
