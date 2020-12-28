package main

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
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
	dsn := "sqlserver://SqlServer_username:SqlServer_Password@SqlServer_IP:SqlServer_port?database=SqlServer_Database"
	db, err = gorm.Open(sqlserver.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Connect database err ,%v\n", err)
	}
	_ = db.AutoMigrate(&UserInfos{})
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
	r.JSON(http.StatusCreated, gin.H{
		"code": 201,
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

func Login(l *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := l.BindJSON(&input)
	if err != nil {
		l.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Input error please check it!",
		})
		return
	}
	//Whether Username exist or not
	var user UserInfos
	db.Where("username = ?", input.Username).First(&user)
	if user.ID == 0 {
		l.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "There are no username like this ,please register",
		})
		return
	}
	//Whether password correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		l.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "Incorrect Password Please check it",
		})
		return
	}
	token, err := ReleaseToken(user)
	if err != nil {
		l.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "server error",
			//log.Printf("token gernerate error : %v",err)
		})
		return
	}
	l.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    gin.H{"token": token},
		"message": "登录成功",
	})
}

var jwtKey = []byte("a_secret_crect")

type Claims struct {
	UserID uint
	jwt.StandardClaims
}

func ReleaseToken(user UserInfos) (string, error) {
	expirationTime := time.Now().Add(7 * 24 * time.Hour)
	claims := &Claims{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "Issure_Code",
			Subject:   "user token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}
