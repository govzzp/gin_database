package main

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	"golang.org/x/crypto/bcrypt"
)

type UserInfos struct {
	gorm.Model
	Telephone string `gorm:"type:char(11);index;unique;not null;" json:"telephone"` //手机号,加索引，唯一，不为空
	Username  string `gorm:"type:varchar(12);" json:"username"`                     //用户名
	Password  string `gorm:"type:varchar(100);"`                                    //密码
	Power     int    `gorm:"type:int;" json:"power"`                                //用户权限 1为正常用户 2为管理员用户
	Sex       int    `gorm:"type:int;default:0;" json:"sex"`                        //用户性别
	Age       int    `gorm:"type:int;" json:"age"`                                  //用户年龄，存储的是时间戳字符串
	State     int    `gorm:"type:int;default:0;" json:"-"`                          //用户状态
}

var db *gorm.DB

func init() {
	if err := InitDB(); err != nil {
		panic(err)
	}
}

func InitDB() (err error) {
	db, err = gorm.Open("mssql", "sqlserver://SqlServer_username:SqlServer_Password@SqlServer_IP:SqlServer_port?database=SqlServer_Database")
	if err != nil {
		fmt.Printf("Connect database err ,%v\n", err)
	}
	db.AutoMigrate(&UserInfos{})
	return
}

//func Getdb() *gorm.DB {
//	return db
//}
func main() {
	r := gin.Default()
	r.POST("/useradd", Register)
	panic(r.Run(":8080"))
}

func Register(r *gin.Context) {
	num := `[0-9]{1}`
	aZ := `[a-z]{1}`
	AZ := `[A-Z]{1}`
	symbol := `[!@#~$%^&*()+|_]{1}`
	var input struct {
		Telephone string `json:"telephone" binding:"required"`
		Username  string `json:"username" binding:"required"`
		Password  string `json:"password" binding:"required"`
		Power     int    `json:"power" binding:"required"`
		Sex       int    `json:"sex" binding:"required"`
		Age       int    `json:"age" binding:"required"`
	}
	err := r.BindJSON(&input)
	if err != nil {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "input error please check your input",
		})
		return
	}
	if len(input.Telephone) != 11 {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Telephone number must be 11",
		})
		return
	}
	if len(input.Password) < 6 && len(input.Password) > 18 {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Password must less than 18 words and more than 6 words",
		})
		return
	}
	if b, err := regexp.MatchString(num, input.Password); !b || err != nil {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "Password must have numbers 0-9!!!",
		})
		return
	}
	if b, err := regexp.MatchString(aZ, input.Password); !b || err != nil {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "Password must have lowercase letters.",
		})
		return
	}
	if b, err := regexp.MatchString(AZ, input.Password); !b || err != nil {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "Password must have Capitals.",
		})
		return
	}
	if b, err := regexp.MatchString(symbol, input.Password); !b || err != nil {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "Password must have less than one Punctuation.",
		})
		return
	}
	if isUsernameExist(db, input.Username) {
		r.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "User already exist please login",
		})
		return
	}
	hasedpassword, erro := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if erro != nil {
		r.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "internet server error",
		})
		return
	}
	newUser := UserInfos{
		Username:  input.Username,
		Password:  string(hasedpassword),
		Telephone: input.Telephone,
		Power:     input.Power,
		Sex:       input.Sex,
		Age:       input.Age,
	}
	db.Create(&newUser)
	r.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Register successful",
	})
}
func isUsernameExist(db *gorm.DB, username string) bool {
	var user UserInfos
	db.Where("username = ?", username).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

//func Login(l *gin.Context) {
//
//}
