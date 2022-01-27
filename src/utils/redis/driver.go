package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"todoList/config"
)

func Connect() (client *redis.Client) {

	cfg, error := config.Config()
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	host := cfg.Section("cache").Key("host").String()
	port := cfg.Section("cache").Key("port").String()
	password := cfg.Section("cache").Key("password").String()

	addr := fmt.Sprintf("%s:%s", host, port)
	client = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	pong, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis server error:", err.Error())
	} else {
		fmt.Println("ping:", pong)
	}

	return
}

func Close(client *redis.Client) {
	client.Close()
}
