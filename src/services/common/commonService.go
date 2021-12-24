package common

import (
	"context"
	"fmt"
	"strconv"
	"time"
	redis "todoList/src/utils/redis"
)

var ctx = context.Background()

func GetUID(name string) (uid int) {

	currentDate := time.Now().Format("20060102")
	key := fmt.Sprintf("%s-%s", name, currentDate)
	id := Incr(key)
	uidStr := fmt.Sprintf("%s%d", currentDate, id)
	uid, _ = strconv.Atoi(uidStr)

	return
}

func Incr(name string) (id int) {
	client := redis.Connect()
	defer redis.Close(client)

	result, _ := client.Incr(ctx, name).Result()
	client.Expire(ctx, name, 24 * 60 * 60 * time.Second)

	id = int(result)
	return
}
