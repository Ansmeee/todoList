package common

import (
	"context"
	"crypto/md5"
	"fmt"
	"strconv"
	redis "todoList/src/utils/redis"
)

var ctx = context.Background()
func GetUID(name string) (uid string) {
	id := Incr(name)
	data := []byte(strconv.Itoa(id))
	uid = fmt.Sprintf("%x", md5.Sum(data))
	return
}

func Incr(name string) (id int)  {
	client := redis.Connect()
	defer redis.Close(client)

	result, _ := client.Incr(ctx, name).Result()
	id = int(result)
	return
}