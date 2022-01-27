package mysql

import (
	"fmt"
	"gopkg.in/ini.v1"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"todoList/config"
)

type mySQLConnection struct {
	Host     string
	Port     uint
	Username string
	Password string
	Database string
}

func generateDSN(config *ini.File) string {
	username := config.Section("database").Key("username")
	password := config.Section("database").Key("password")
	host := config.Section("database").Key("host")
	port := config.Section("database").Key("port")
	dbname := config.Section("database").Key("database")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	return dsn
}

func Connect(dbname string) (db *gorm.DB, error error) {
	cfg, error := config.Config()
	if error != nil {
		fmt.Println(error.Error())
		return
	}

	dsn := generateDSN(cfg)
	db, error = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	return
}
