package model

import (
	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	BookCode   string `gorm:"type:varchar(100);index;unique;not null" json:"bookcode"`
	ISBN       string `gorm:"type:varchar(50);not null;" json:"isbn"`
	BookName   string `gorm:"type:varchar(200);not null;" json:"bookname"`
	Author     string `gorm:"type:varchar(100);not null;" json:"author"`
	Press      string `gorm:"type:varchar(100);not null;" json:"press"`    //出版社
	Category   string `gorm:"type:varchar(100);not null;" json:"category"` //种类
	GetBooknum string `gorm:"type:varchar(200);not null;" json:"getbooknum"`
	Position   string `gorm:"type:varchar(200);not null;" json:"position"`
}
type Borrow struct {
	gorm.Model
	Username string `gorm:"type:varchar(12);not null;" json:"username"`
	ISBN     string `gorm:"type:varchar(50);not null;" json:"isbn"`
	BookName string `gorm:"type:varchar(200);not null;"  json:"bookname"`
	BookCode string `gorm:"type:varchar(100);not null;unique;index;" json:"bookcode"`
}

type UserInfos struct {
	gorm.Model
	Telephone string  `gorm:"type:char(11);index;unique;not null;" json:"telephone"` //手机号,加索引，唯一，不为空
	Username  string  `gorm:"type:varchar(12);" json:"username"`                     //用户名一般以学号或者教师工号替代但是会加入admin和其他用户名
	Password  string  `gorm:"type:varchar(100);"`                                    //密码
	Power     int     `gorm:"type:int;" json:"power"`                                //用户权限 1为正常用户 2为管理员用户
	Sex       int     `gorm:"type:int;default:0;" json:"sex"`                        //用户性别
	Age       int     `gorm:"type:int;" json:"age"`                                  //用户年龄
	Money     float32 `gorm:"type:float(10,4); default:0.0;"  json:"money"`          //用户有没有欠款
}
type Token struct {
	gorm.Model
	UserID uint
	Token  string
}
