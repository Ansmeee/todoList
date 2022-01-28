package bootstrap

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"os"
	"strings"
	"todoList/config"
	"todoList/src/controllers/middleware"
	"todoList/src/routers"
)

func StartEngine() {
	args := os.Args

	var error error
	var cfg *ini.File

	if len(args) == 2 {
		cfgFile := args[1]
		cfg, error = ini.Load(cfgFile)
	}

	if error != nil || cfg == nil {
		cfg, error = ini.Load("config.ini")
	}

	// 初始化配置文件
	config.InitConfig(cfg)

	engine := gin.Default()
	routerGroup := new(routers.RouterGroup)
	var routers = engine.Group("rest").Use(middleware.Auth)
	routerGroup.InitRouter(routers)

	ip := strings.TrimSpace(cfg.Section("server").Key("ip").String())
	if len(ip) == 0 {
		ip = "127.0.0.1"
	}

	port := strings.TrimSpace(cfg.Section("server").Key("port").String())
	if len(port) == 0 {
		port = "8000"
	}

	addr := fmt.Sprintf("%s:%s", ip, port)
	engine.Run(addr)
}
