package Controller

import (
	"github.com/gin-gonic/gin"
	"go_free/common"
	"go_free/dto"
	"go_free/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
	"regexp"
)

func Register(r *gin.Context) {
	db := common.Getdb()
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
	newUser := model.UserInfos{
		Username:  input.Username,
		Password:  string(hasedpassword),
		Telephone: input.Telephone,
		Power:     input.Power,
		Sex:       input.Sex,
		Age:       input.Age,
	}
	tx := db.Begin()
	if tx.Create(&newUser).RowsAffected != 1  {
		r.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Internet Server Error!",
		})
		tx.Commit()
		return
	}
	r.JSON(http.StatusCreated, gin.H{
		"code": 201,
		"msg":  "Register successful",
	})
}
func isUsernameExist(db *gorm.DB, username string) bool {
	var user model.UserInfos
	db.Where("username = ?", username).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

func Login(l *gin.Context) {
	db := common.Getdb()
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
	var user model.UserInfos
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
	token, err := common.ReleaseToken(user)
	if err != nil {
		l.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "server error",
			//log.Printf("token gernerate error : %v",err)
		})
		return
	}
	tx := db.Begin()
	if tx.Create(&model.Token{UserID: user.ID, Token: token}).RowsAffected != 1 {
		l.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Internet Server Error!",
		})
		tx.Commit()
		return
	}
	l.JSON(http.StatusOK, gin.H{
		"code":    200,
		"data":    gin.H{"token": token},
		"message": "登录成功",
	})
}
func Info(c *gin.Context) {
	user, _ := c.Get("user")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": dto.ToUserDto(user.(model.UserInfos))},
	})
}
