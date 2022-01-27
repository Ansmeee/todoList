package config

import (
	"errors"
	"gopkg.in/ini.v1"
)

var config *ini.File
func InitConfig(cfg *ini.File)  {
	config = cfg
}

func Config() (*ini.File, error) {
	var error error

	if config == nil {
		error = errors.New("load config file failed")
	}

	return config, error
}