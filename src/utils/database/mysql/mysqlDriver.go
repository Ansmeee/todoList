package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type mySQLConnection struct {
	Host string
	Port uint
	Username string
	Password string
	Database string
}

func Connect(dbname string) (db *gorm.DB, error error)  {
	dsn := fmt.Sprintf("dev:dev007@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True&loc=Local")
	db, error = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return
}
