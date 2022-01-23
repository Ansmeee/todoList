package user

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"todoList/src/config"
	"todoList/src/models/user"
	"todoList/src/utils/database"
	"todoList/src/utils/redis"
)

type UserService struct{}

var thisService = &UserService{}
var thisModel = &user.UserModel{}
var ctx = context.Background()

func (service *UserService) FindeByEmail(email string) (error error, data *user.UserModel) {
	client := redis.Connect()
	defer redis.Close(client)

	email = strings.TrimSpace(email)
	userCacheKey := fmt.Sprintf("user:%s", email)
	cacheData, err := client.Get(ctx, userCacheKey).Bytes()
	if err != nil {
		fmt.Println(err.Error())
	}

	json.Unmarshal(cacheData, data)
	if data != nil {
		return
	}

	db := database.Connect("")
	defer database.Close(db)

	err = db.Model(thisModel).Where("email = ?", email).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		error = errors.New("系统异常")
		return
	}

	err = rebuildCacke(userCacheKey, *data)
	if err != nil {
		fmt.Println("缓存更新失败")
	}

	return
}

func (service *UserService) FindByID(id int) (error error, data *user.UserModel) {
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
	Account string `form:"account"`
	Auth    string `form:"auth"`
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

func (service *UserService) SignIn(data *SigninForm) (token string, error error) {
	err, userInfo := thisService.FindeByEmail(data.Account)
	if err != nil {
		 error = errors.New("用户信息异常")
		 return
	}

	fmt.Println(userInfo)

	if userInfo.Id == 0 {
		error = errors.New("用户不存在")
		return
	}

	if ! userInfo.Active() {
		error = errors.New("该用户已删除")
		return
	}

	token, err = thisService.GenerateToken(userInfo)
	if err != nil {
		error = errors.New("请重试")
		return
	}

	res := thisService.LoginByToken(token, *userInfo)
	if res != true {
		error = errors.New("请重试")
		return
	}

	return
}

func (UserService) LogoutByToken(token string) (error error)  {
	client := redis.Connect()
	defer redis.Close(client)

	error = client.Del(ctx, token).Err()
	return
}

func (UserService) LoginByToken(token string, data user.UserModel) bool  {
	client := redis.Connect()
	defer redis.Close(client)

	encodeData, error := json.Marshal(data.Id)
	if error != nil {
		fmt.Println(error.Error())
		return false
	}

	expireTime := 24 * 60 * 60 * time.Second
	if _, error := client.Set(ctx, token, encodeData, expireTime).Result(); error != nil {
		fmt.Println(error.Error())
		return false
	}

	return true
}

type SignupForm struct {
	Email    string `form:"email"`
	PassWord string `form:"password"`
	Auth     string `form:"auth"`
}
func (service *UserService) SignUp(data *SignupForm) (err error) {
	db := database.Connect("")
	defer database.Close(db)

	err = db.Transaction(func(tx *gorm.DB) error {
		var newUser user.UserModel
		newUser.Email = data.Email
		if tx.Model(thisModel).Create(&newUser).Error != nil {
			return errors.New("用户信息存储失败")
		}

		var userAuth user.AuthModel
		userAuth.Email = newUser.Email
		userAuth.Auth  = data.Auth
		if tx.Model(&user.AuthModel{}).Create(&userAuth).Error != nil {
			return errors.New("用户信息存储失败")
		}

		return nil
	})

	return
}

type AttrForm struct {
	Id    int    `form:"id"`
	Key   string `form:"key"`
	Value string `form:"value"`
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

	if err := rebuildCacke(userCacheKey, user); err != nil {
		error = errors.New("缓存更新失败")
	}

	return
}

func rebuildCacke(cacheKey string, data interface{}) (error error) {
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

	if total == 0 { return }

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

func (UserService) GenerateToken(userInfo *user.UserModel) (token string, error error)  {
	header := map[string]string{"typ": "JWT", "alg": "HS256"}
	headerByte, _ := json.Marshal(header)
	encodingHeader := base64.StdEncoding.EncodeToString(headerByte)

	payload := map[string]interface{}{"account": userInfo.Id, "name": userInfo.Name, "expiredat": time.Now().Add(24 * time.Hour)}
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

func (UserService) GetUserInfoByToken(token string) (data *user.UserModel, error error)  {
	client := redis.Connect()
	defer redis.Close(client)

	cacheData, err := client.Get(ctx, token).Bytes()

	if err != nil {
		error = errors.New("用户信息获取失败")
		return
	}

	var email string
	err = json.Unmarshal(cacheData, &email)
	if err != nil {
		error = errors.New("用户信息获取失败")
		return
	}

	if len(email) == 0 {
		error = errors.New("用户信息获取失败")
		return
	}

	err, data = thisService.FindeByEmail(email)
	return
}