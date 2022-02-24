package config

import (
	"errors"
	"gopkg.in/ini.v1"
	"os"
)

var config *ini.File
func InitConfig()  {
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

	config = cfg
}

func Config() (*ini.File, error) {
	var error error

	if config == nil {
		error = errors.New("load config file failed")
	}

	return config, error
}