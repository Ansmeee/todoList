package common

import (
	"context"
	"crypto/md5"
	"fmt"
	redis "todoList/src/utils/redis"
)

var ctx = context.Background()
func GetUID() (uid string, error error) {
	client := redis.Connect()
	defer redis.Close(client)

	id := client.Incr(ctx, "uid").String()
	fmt.Println(id)

	if error != nil {
		return
	}

	data := []byte(id)
	uid = fmt.Sprintf("%x", md5.Sum(data))
	return
}