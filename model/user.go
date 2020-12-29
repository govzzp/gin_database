package model

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	ISBN       string `gorm:"type:varchar(50);index;unique;not null;" json:"isbn"`
	bookName   string `gorm:"type:varchar(200);" json:"bookname"`
	Author     string `gorm:"type:varchar(100);" json:"author"`
	Press      string `gorm:"type:varchar(100);" json:"press"`
	Category   string `gorm:"type:varchar(100);" json:"category"`
	getBooknum string `gorm:"type:varchar(200);" json:"getbooknum"`
}
type Borrow struct {
	gorm.Model
	Username   string `gorm:"type:varchar(12);index;unique;not null;" json:"username"`
	ISBN       string `gorm:"type:varchar(50)；not null;" json:"isbn"`
	bookName   string `gorm:"type:varchar(200);"  json:"book_name"`
	getBookNum string `gorm:"type:varchar(200);"json:"getbooknum"`
}

type UserInfos struct {
	gorm.Model
	Telephone string `gorm:"type:char(11);index;unique;not null;" json:"telephone"` //手机号,加索引，唯一，不为空
	Username  string `gorm:"type:varchar(12);" json:"username"`                     //用户名一般以学号或者教师工号替代但是会加入admin和其他用户名
	Password  string `gorm:"type:varchar(100);"`                                    //密码
	Power     int    `gorm:"type:int;" json:"power"`                                //用户权限 1为正常用户 2为管理员用户
	Sex       int    `gorm:"type:int;default:0;" json:"sex"`                        //用户性别
	Age       int    `gorm:"type:int;" json:"age"`                                  //用户年龄
}
