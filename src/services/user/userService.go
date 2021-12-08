package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	"todoList/src/models/user"
	"todoList/src/utils/database"
	"todoList/src/utils/redis"
)

type UserService struct{}

var thisService = &UserService{}
var thisModel = &user.UserModel{}
var ctx = context.Background()

func (service *UserService) FindeByEmail(email string) (error error, data user.UserModel) {
	client := redis.Connect()
	defer redis.Close(client)

	email = strings.TrimSpace(email)
	userCacheKey := fmt.Sprintf("user:%s", email)
	cacheData, err := client.Get(ctx, userCacheKey).Bytes()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = json.Unmarshal(cacheData, &data)
	if err == nil {
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

	err = rebuildCacke(userCacheKey, data)
	if err != nil {
		fmt.Println("缓存更新失败")
	}

	return
}

func (service *UserService) FindByID(id string) (error error, data user.UserModel) {
	client := redis.Connect()
	defer redis.Close(client)

	userCacheKey := fmt.Sprintf("user:%d", id)
	cacheData, err := client.Get(ctx, userCacheKey).Bytes()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = json.Unmarshal(cacheData, &data)
	if err == nil {
		return
	}

	db := database.Connect("")
	defer database.Close(db)

	err = db.Model(thisModel).First(&data, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			error = errors.New("不存在该记录")
			return
		}

		error = errors.New("获取失败")
		return
	}

	err = rebuildCacke(userCacheKey, data)
	if err != nil {
		fmt.Println("缓存更新失败")
	}

	return
}

type SigninForm struct {
	Account string `form:"account"`
	Auth    string `form:"auth"`
}

func (service *UserService) SignIn(data *SigninForm) (token string, error error) {
	err, userInfo := thisService.FindeByEmail(data.Account)
	if err != nil {
		 error = errors.New("用户信息异常")
		 return
	}

	if userInfo.Id == "" {
		error = errors.New("用户信息异常")
		return
	}

	if ! userInfo.OnJob() {
		error = errors.New("该用户已删除")
		return
	}

	token, err = thisService.GenerateToken(&userInfo)
	if err != nil {
		error = errors.New("用户信息异常")
		return
	}

	res := thisService.LoginByToken(token, userInfo)
	if res != true {
		error = errors.New("登陆失败，请重试")
		return
	}

	return
}

func (UserService) LoginByToken(token string, data user.UserModel) bool  {
	client := redis.Connect()
	defer redis.Close(client)

	encodeData, error := json.Marshal(data.Email)
	if error != nil {
		fmt.Println(error.Error())
		return false
	}

	expireTime := time.Second * 60 * 60
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
	token = time.Now().String()
	// TODO
	return
}

func (UserService) GetUserInfoByToken(token string) (data user.UserModel, error error)  {
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