package common

import (
	"fmt"
	"github.com/spf13/viper"
	"go_free/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
)

var db *gorm.DB

func init() {
	if err := InitDB(); err != nil {
		panic(err)
	}
}
func InitConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "\\config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}
func InitDB() (err error) {
	InitConfig()
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	database := viper.GetString("datasource.database")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
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
