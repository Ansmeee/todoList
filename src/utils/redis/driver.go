package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func Connect() (client *redis.Client){
	client = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		Password: "ansme@redis",
		DB: 0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis server error:", err.Error())
	} else {
		fmt.Println("ping:", pong)
	}

	return
}

func Close(client *redis.Client)  {
	client.Close()
}