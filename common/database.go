package common

import (
	"fmt"
	"go_free/model"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	if err := InitDB(); err != nil {
		panic(err)
	}
}

func InitDB() (err error) {
	dsn := "sqlserver://SqlServer_username:SqlServer_Password@SqlServer_IP:SqlServer_port?database=SqlServer_Database"
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Connect database err ,%v\n", err)
	}
	_ = db.AutoMigrate(&model.UserInfos{})
	_ = db.AutoMigrate(&model.Book{})
	_ = db.AutoMigrate(&model.Borrow{})
	return
}

func Getdb() *gorm.DB {
	return db
}
