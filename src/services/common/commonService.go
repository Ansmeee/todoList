package common

import (
	"context"
	"fmt"
	"github.com/bwmarrin/snowflake"
	"strconv"
	"time"
	"todoList/config"
	"todoList/src/utils/redis"
)

var ctx = context.Background()

func GetUID() (uid int64) {

	cfg, _ := config.Config()
	nodeID := cfg.Section("environment").Key("app_node").String()
	nodeNo, err := strconv.ParseInt(nodeID, 10, 64)
	if err != nil {
		fmt.Println("GetUID Error: ", err.Error())
		return
	}

	node, err := snowflake.NewNode(nodeNo)
	if err != nil {
		fmt.Println("GetUID Error: ", err.Error())
		return
	}

	uid = node.Generate().Int64()
	return
}

func Incr(name string) (id int) {
	client := redis.Connect()
	defer redis.Close(client)

	result, _ := client.Incr(ctx, name).Result()
	client.Expire(ctx, name, 24*60*60*time.Second)

	id = int(result)
	return
}
