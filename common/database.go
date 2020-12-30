package common

import (
	"fmt"
	"go_free/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	if err := InitDB(); err != nil {
		panic(err)
	}
}

func InitDB() (err error) {
	dsn := "mysql_username:mysql_password@tcp(mysql_IP:mysql_port)/mysql_database?charset=utf8mb4&parseTime=True&loc=Local"
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Connect database err ,%v\n", err)
	}
	_ = db.AutoMigrate(&model.UserInfos{})
	_ = db.AutoMigrate(&model.Book{})
	_ = db.AutoMigrate(&model.Borrow{})
	_ = db.AutoMigrate(&model.Token{})
	return
}

func Getdb() *gorm.DB {
	return db
}
