package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"gorm.io/gorm"
	"log"
	"math/rand"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
	cfg "todoList/config"
	"todoList/src/config"
	"todoList/src/models/user"
	"todoList/src/services/common"
	"todoList/src/services/mailSVC"
	"todoList/src/services/smsSVC"
	"todoList/src/utils/database"
	"todoList/src/utils/redis"
)

type UserService struct{}

var thisService = &UserService{}
var thisModel = &user.UserModel{}
var ctx = context.Background()

func (s *UserService) FindByPhone(phone string) (*user.UserModel, error) {
	db := database.Connect("")
	defer database.Close(db)

	data := new(user.UserModel)
	err := db.Model(thisModel).Where("phone = ?", phone).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		log.Println("mysql query error:", err)
		return nil, errors.New("系统异常")
	}

	return data, nil
}

func (service *UserService) FindByEmail(email string) (error error, data *user.UserModel) {
	db := database.Connect("")
	defer database.Close(db)

	err := db.Model(thisModel).Where("email = ?", email).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		error = errors.New("系统异常")
		return
	}

	return
}

func (service *UserService) FindByID(id string) (error error, data *user.UserModel) {
	db := database.Connect("")
	defer database.Close(db)

	data = new(user.UserModel)
	error = db.Where("uid = ?", id).Find(data).Error
	if error != nil {
		if errors.Is(error, gorm.ErrRecordNotFound) {
			error = errors.New("该用户不存在")
			return
		}

		error = errors.New("获取失败")
		return
	}

	return
}

type SigninForm struct {
	Account  string `form:"account"`
	Password string `form:"password"`
	Way      string `form:"way"`
}

func (service *UserService) SignOut(token string) (error error) {
	if len(token) == 0 {
		return
	}

	error = thisService.LogoutByToken(token)
	if error != nil {
		error = errors.New("系统异常")
	}

	return
}

func signinByAccount(data *SigninForm) (string, error) {
	var err error
	var userInfo *user.UserModel
	if strings.Contains(data.Account, "@") {
		err, userInfo = thisService.FindByEmail(data.Account)
	} else {
		userInfo, err = thisService.FindByPhone(data.Account)
	}

	if err != nil || userInfo == nil || userInfo.Id == "" {
		return "", errors.New("该用户不存在")
	}

	if err = thisService.AuthPassword(userInfo.Id, data.Password); err != nil {
		return "", errors.New("用户名或密码不正确")
	}

	token, err := thisService.GenerateToken(userInfo)
	if err != nil {
		return "", errors.New("请重试")
	}

	res := thisService.LoginByToken(token, *userInfo)
	if res != true {
		return "", errors.New("请重试")
	}

	return token, nil
}

func signinBySMS(data *SigninForm) (string, error) {
	code := thisService.SMSCode(data.Account)
	if code == "" || code != data.Password {
		return "", errors.New("请输入正确的验证码")
	}

	userInfo, err := thisService.FindByPhone(data.Account)
	if err != nil {
		return "", errors.New("用户信息异常")
	}

	if err == nil && userInfo == nil {
		signupForm := new(SignupForm)
		signupForm.Phone = data.Account
		userInfo, err = signUpWithData(signupForm)
		if err != nil {
			return "", errors.New("用户信息初始化失败")
		}
	}

	token, err := thisService.GenerateToken(userInfo)
	if err != nil {
		return "", errors.New("请重试")
	}

	res := thisService.LoginByToken(token, *userInfo)
	if res != true {
		return "", errors.New("请重试")
	}

	return token, nil
}

func (service *UserService) SignIn(data *SigninForm) (string, error) {

	if data.Way == "sms" {
		return signinBySMS(data)
	}

	return signinByAccount(data)
}

func (UserService) AuthPassword(account, password string) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	auth := new(user.AuthModel)
	error = db.Where("account = ? and auth = ?", account, password).First(auth).Error
	return
}

func (UserService) LogoutByToken(token string) (error error) {
	client := redis.Connect()
	defer redis.Close(client)

	error = client.Del(ctx, token).Err()
	return
}

func (UserService) LoginByToken(token string, data user.UserModel) bool {
	client := redis.Connect()
	defer redis.Close(client)

	encodeData, error := json.Marshal(data.Id)
	if error != nil {
		fmt.Println(error.Error())
		return false
	}

	conf, error := cfg.Config()
	if error != nil {
		fmt.Println(error.Error())
		return false
	}

	tokenLifeTime := conf.Section("cache").Key("token_life_time").MustInt(86400)
	expireTime := time.Second * time.Duration(tokenLifeTime)
	if _, error := client.Set(ctx, token, encodeData, expireTime).Result(); error != nil {
		fmt.Println(error.Error())
		return false
	}

	return true
}

type SignupForm struct {
	Email    string `form:"account"`
	Phone    string `form:"phone"`
	PassWord string `form:"password"`
	Auth     string `form:"auth"`
}

func signUpWithData(form *SignupForm) (data *user.UserModel, err error) {
	db := database.Connect("")
	defer database.Close(db)

	err = db.Transaction(func(tx *gorm.DB) error {
		newUser := new(user.UserModel)
		newUser.Id = common.GetUID()

		if form.Email != "" {
			newUser.Email = form.Email
			newUser.Name = form.Email
		}

		if form.Phone != "" {
			newUser.Phone = form.Phone
			newUser.Name = form.Phone
		}

		if tx.Model(thisModel).Create(&newUser).Error != nil {
			return errors.New("注册失败")
		}

		if form.Email != "" {
			var userAuth user.AuthModel
			userAuth.Account = newUser.Id
			userAuth.Auth = form.Auth
			if tx.Model(&user.AuthModel{}).Create(&userAuth).Error != nil {
				return errors.New("注册失败")
			}
		}

		data = newUser
		return nil
	})

	return
}

//func (service *UserService) SignUp(data *SignupForm) (user *user.UserModel, error error) {
//	error = errors.New("注册失败")
//	return
//}

type ResetPassForm struct {
	Password string `form:"password"`
	Auth     string `form:"auth"`
	Account  string `form:"account"`
}

func (UserService) ResetPass(form *ResetPassForm) (error error) {
	if user.User().Id != form.Account {
		error = errors.New("用户信息异常")
		return
	}

	if form.Password != form.Auth {
		error = errors.New("两次输入的密码不一致")
		return
	}

	db := database.Connect("")
	defer database.Close(db)

	userauth := new(user.AuthModel)

	if err := db.Model(userauth).Where("account = ?", user.User().Id).First(userauth).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			userauth.Account = user.User().Id
			userauth.Auth = form.Auth

			if db.Model(&user.AuthModel{}).Create(&userauth).Error != nil {
				return errors.New("密码更新失败")
			}

			return
		}

		log.Println("update password error:", err)
		error = errors.New("用户信息异常")
		return
	}

	fmt.Println("here")
	if form.Password == userauth.Auth {
		error = errors.New("新密码与旧密码不能相同")
		return
	}

	if err := db.Model(&userauth).Where("account = ?", user.User().Id).Update("auth", form.Password).Error; err != nil {
		error = errors.New("密码重置失败")
		return
	}

	return
}

type AttrForm struct {
	Id    string `form:"id"`
	Key   string `form:"key"`
	Value string `form:"value"`
	Code  string `form:"code"`
}

func (UserService) UpdateAttr(user *user.UserModel, key string, value interface{}) (error error) {
	db := database.Connect("")
	defer database.Close(db)

	updateData := map[string]interface{}{"updated_at": time.Now().Format("2006-01-02 15:01:05")}
	updateData[key] = value

	error = db.Model(user).Where("uid = ?", user.Id).Updates(updateData).Error
	return
}

func (service *UserService) Update(user, data *user.UserModel) (error error) {
	client := redis.Connect()
	defer redis.Close(client)

	userCacheKey := fmt.Sprintf("user:%d", data.Id)
	if _, err := client.Del(ctx, userCacheKey).Result(); err != nil {
		fmt.Println(err.Error())
	}

	db := database.Connect("")
	defer database.Close(db)

	if db.Model(&user).Updates(data).Error != nil {
		error = errors.New("系统异常")
	}

	if err := rebuildCache(userCacheKey, user); err != nil {
		error = errors.New("缓存更新失败")
	}

	return
}

func rebuildCache(cacheKey string, data interface{}) (error error) {
	client := redis.Connect()
	defer redis.Close(client)

	encodeData, error := json.Marshal(data)
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	expiratedAt := time.Second * 60 * 60
	if _, error := client.Set(ctx, cacheKey, encodeData, expiratedAt).Result(); error != nil {
		fmt.Println(error.Error())
	}

	return
}

type QueryParams struct {
	Keywords string `json:"keywords" form:"keywords"`
	Page     int    `json:"page" form:"page"`
	PageSize int    `json:"page_size" form:"page_size"`
}

func (service *UserService) List(params *QueryParams) (error error, data interface{}, total int64) {
	db := database.Connect("")
	defer database.Close(db)

	data = []user.UserModel{}

	wheres := make([]string, 0)
	if len(params.Keywords) > 0 {
		wheres = append(wheres, fmt.Sprintf("(`name` like '%s%%' or `email` like '%s%%')", params.Keywords, params.Keywords))
	}

	if len(wheres) > 0 {
		if db.Model(thisModel).Where(strings.Join(wheres, " and ")).Count(&total).Error != nil {
			error = errors.New("系统异常")
			return
		}
	} else {
		if db.Model(thisModel).Count(&total).Error != nil {
			error = errors.New("系统异常")
			return
		}
	}

	if total == 0 {
		return
	}

	var userList []user.UserModel
	page := (params.Page - 1) * params.PageSize
	if db.Model(thisModel).Limit(params.PageSize).Offset(page).Find(&userList).Error != nil {
		error = errors.New("系统异常")
		return
	}

	return error, userList, total
}

func (UserService) Delete(user *user.UserModel) (error error) {
	client := redis.Connect()
	defer redis.Close(client)

	userCacheKey := fmt.Sprintf("user:%d", user.Id)
	if _, err := client.Del(ctx, userCacheKey).Result(); err != nil {
		error = errors.New("缓存删除失败")
		return
	}

	db := database.Connect("")
	defer database.Close(db)
	error = db.Delete(&user).Error
	return
}

func (UserService) GenerateToken(userInfo *user.UserModel) (token string, error error) {
	conf, error := cfg.Config()
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	tokenLifeTime := conf.Section("cache").Key("token_life_time").MustInt(86400)
	return generate(userInfo, tokenLifeTime)
}

func generate(userInfo *user.UserModel, tokenLifeTime int) (token string, error error) {
	header := map[string]string{"typ": "JWT", "alg": "HS256"}
	headerByte, _ := json.Marshal(header)
	encodingHeader := base64.StdEncoding.EncodeToString(headerByte)

	expireTime := time.Hour * time.Duration(tokenLifeTime/60/60)
	payload := map[string]interface{}{
		"account":   userInfo.Id,
		"name":      userInfo.Name,
		"expiredat": time.Now().Add(expireTime),
		"icon":      userInfo.Icon,
	}
	payloadByte, _ := json.Marshal(payload)
	encodingPayload := base64.StdEncoding.EncodeToString(payloadByte)
	secret := []byte(config.Secret)

	encodingString := encodingHeader + "." + encodingPayload

	hash := hmac.New(sha256.New, secret)
	hash.Write([]byte(encodingString))

	signature := strings.TrimRight(base64.URLEncoding.EncodeToString(hash.Sum(nil)), "=")
	token = encodingString + "." + signature
	return
}

func (UserService) GetUserInfoByToken(token string) (data *user.UserModel, error error) {
	client := redis.Connect()
	defer redis.Close(client)

	cacheData, err := client.Get(ctx, token).Bytes()
	if err != nil {
		error = errors.New("用户信息获取失败")
		return
	}

	var account string
	err = json.Unmarshal(cacheData, &account)
	if err != nil {
		error = errors.New("用户信息获取失败")
		return
	}

	if account == "" {
		error = errors.New("用户信息获取失败")
		return
	}

	err, data = thisService.FindByID(account)
	return
}

func (UserService) GenerateLocalIconPath() (url string, error error) {
	conf, error := cfg.Config()
	if error != nil {
		return
	}

	iconPath := conf.Section("environment").Key("icon_path").String()
	if "" == iconPath {
		fmt.Println("用户头像文件路径配置异常")
		error = errors.New("上传失败")
		return
	}

	_, error = os.Stat(iconPath)
	if error != nil {
		if os.IsNotExist(error) {

			error = os.MkdirAll(iconPath, os.ModePerm)

			if error != nil {
				fmt.Println("用户头像文件路径创建失败")
				error = errors.New("上传失败")
				return
			}
		} else {
			if error != nil {
				fmt.Println("用户头像文件路径创建失败")
				error = errors.New("上传失败")
				return
			}
		}
	}

	url = iconPath
	return
}

func (*UserService) Verifing(userInfo *user.UserModel) bool {
	client := redis.Connect()
	defer redis.Close(client)
	verifyKey := fmt.Sprintf("email:%s", userInfo.Email)
	if res, err := client.Exists(ctx, verifyKey).Result(); res > 0 && err == nil {
		return true
	}

	return false
}

func (*UserService) Verified(token string) bool {
	client := redis.Connect()
	defer redis.Close(client)

	data, err := client.Get(ctx, token).Bytes()
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	var email string
	err = json.Unmarshal(data, &email)

	if err != nil || email == "" {
		fmt.Println(err.Error())
		return false
	}

	err, u := thisService.FindByEmail(email)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if u.Id == "" {
		return false
	}

	if u.Verified == user.EMAIL_VERIFIED {
		return true
	}

	verifyKey := fmt.Sprintf("email:%s", u.Email)
	err = client.Del(ctx, verifyKey).Err()
	if err != nil {
		fmt.Println("del verifyKey err:", err.Error())
		return false
	}

	db := database.Connect("")
	defer database.Close(db)

	updateData := map[string]interface{}{
		"verified":   user.EMAIL_VERIFIED,
		"updated_at": time.Now().Format("2006-01-02 15:01:05"),
	}

	if db.Model(&user.UserModel{}).Where("uid = ?", u.Id).Updates(updateData).Error == nil {
		return true
	}

	fmt.Println("updateData err:", err.Error())

	return false
}

func (*UserService) VerifyEmail(userInfo *user.UserModel) error {
	if userInfo.Email == "" {
		return errors.New("邮箱为空")
	}

	conf, error := cfg.Config()
	if error != nil {
		fmt.Println(error.Error())
		return errors.New("读取配置失败")
	}

	tokenLifeTime := conf.Section("cache").Key("email_token_life_time").MustInt(1200)
	host := conf.Section("environment").Key("app_host").String()
	token, error := generate(userInfo, tokenLifeTime)
	if error != nil {
		fmt.Println(error.Error())
		return errors.New("Token 生成失败")
	}

	client := redis.Connect()
	defer redis.Close(client)

	encodeData, error := json.Marshal(userInfo.Email)
	if error != nil {
		fmt.Println(error.Error())
		return errors.New("Token 存储失败")
	}

	expireTime := time.Second * time.Duration(tokenLifeTime)
	if _, error := client.Set(ctx, token, encodeData, expireTime).Result(); error != nil {
		fmt.Println(error.Error())
		return errors.New("Token 存储失败")
	}

	verifyKey := fmt.Sprintf("email:%s", userInfo.Email)
	if _, error := client.Set(ctx, verifyKey, "", expireTime).Result(); error != nil {
		fmt.Println(error.Error())
		return errors.New("验证信息存储失败")
	}

	subject := "土豆清单（ToDoo）邮箱验证"
	url := fmt.Sprintf("%s/email/verify?token=%s", host, token)
	content := fmt.Sprintf("请点击此链接进行邮箱验证，20 分钟内有效: %s", url)

	mailSVC := new(mailSVC.MailSVC)
	if error := mailSVC.SendText(subject, content, userInfo.Email); error != nil {
		fmt.Println(error.Error())
		return errors.New("邮件发送失败")
	}

	return nil
}

func (UserService) SaveIcon2QN(file multipart.File, fileHeader *multipart.FileHeader) (url string, error error) {
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
		error = errors.New("上传失败")
		return
	}

	url = storage.MakePublicURL(imgHost, ret.Key)
	return
}

func (*UserService) SMSCode(account string) string {
	client := redis.Connect()
	defer redis.Close(client)

	key := fmt.Sprintf("smscode:%s", account)
	code, err := client.Get(ctx, key).Result()
	if err != nil {
		return ""
	}

	return code
}

func (*UserService) Frequently(account string) bool {
	client := redis.Connect()
	defer redis.Close(client)

	key := fmt.Sprintf("sendsms:num:%s", account)
	num, err := client.Get(ctx, key).Int()
	if err != nil {
		return false
	}

	return num >= 5
}

func (*UserService) SendSMS(account string) bool {
	if account == "" {
		return false
	}

	code := generateSMSCode(account)
	if code == "" {
		log.Println("SendSMS Error: generate code failed")
		return false
	}
	conf, error := cfg.Config()
	if error != nil {
		fmt.Println(error.Error())
		return false
	}

	env := conf.Section("environment").Key("app_mode").String()
	if env == "production" {
		sms, err := smsSVC.NewSMSSVC()
		if err != nil {
			log.Println("SendSMS Error:", err)
			return false
		}

		err = sms.SendCode(code, account)
		if err != nil {
			return false
		}
	}

	client := redis.Connect()
	defer redis.Close(client)
	key := fmt.Sprintf("sendsms:num:%s", account)
	if err := client.Incr(ctx, key).Err(); err != nil {
		log.Println("sendsms num incr error:", err)
	}

	log.Println("send sms code ok:", code)
	return true
}

func generateSMSCode(account string) string {
	client := redis.Connect()
	defer redis.Close(client)

	key := fmt.Sprintf("smscode:%s", account)

	rand.Seed(time.Now().Unix())
	codeS := make([]string, 4)
	for idx := 0; idx < 4; idx++ {
		codeS = append(codeS, strconv.Itoa(rand.Intn(10)))
	}

	code := strings.Join(codeS, "")

	if code == "" {
		return code
	}

	if err := client.Set(ctx, key, code, 10*60*time.Second).Err(); err != nil {
		log.Println("code save err:", err)
		return ""
	}

	return code
}
