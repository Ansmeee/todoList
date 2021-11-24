package database

import (
	"todoList/src/utils/database/mysql"
	"fmt"
	"gorm.io/gorm"
)

type Driver struct {
}

func Connect(name string) (db *gorm.DB){
	var err error
	db, err = mysql.Connect(name)
	if err != nil {
		fmt.Println("database connect error:", err.Error())
	}

	return
}

func Close(db *gorm.DB)  {
	database, _ := db.DB()
	database.Close()
}