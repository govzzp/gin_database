package Controller

import (
	"go_free/common"
	"go_free/dto"
	"go_free/model"
	"go_free/response"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"input error please check your input",
		)
		return
	}
	if len(input.Telephone) != 11 {
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Telephone number must be 11 words!",
		)
		return
	}
	if len(input.Password) < 6 && len(input.Password) > 18 {
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Password must less than 18 words and more than 6 words",
		)
		return
	}
	if b, err := regexp.MatchString(num, input.Password); !b || err != nil {
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Password must have numbers 0-9!!!",
		)
		return
	}
	if b, err := regexp.MatchString(aZ, input.Password); !b || err != nil {
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Password must have lowercase letters.",
			)
		return
	}
	if b, err := regexp.MatchString(AZ, input.Password); !b || err != nil {
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Password must have Capitals.",
		)
		return
	}
	if b, err := regexp.MatchString(symbol, input.Password); !b || err != nil {
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Password must have less than one Punctuation.",
			)
		return
	}
	if isUsernameExist(db, input.Username) {
		response.Response(r,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"User already exist please login",
			)
		return
	}
	hasedpassword, erro := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if erro != nil {
		response.Response(r,
			http.StatusInternalServerError,
			500,
			nil,
			"internet server error",
			)
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
	if tx.Create(&newUser).RowsAffected != 1 {
		response.Response(r,
			http.StatusInternalServerError,
			500,
			nil,
			"Internet Server Error!",
			)
		tx.Commit()
		return
	}
	response.Success(r,
		nil,
		"Register Successful",
		)
}
func isUsernameExist(db *gorm.DB, username string) bool {
	var user model.UserInfos
	db.Where("username = ?", username).First(&user)
	return user.ID != 0
}

func Login(l *gin.Context) {
	db := common.Getdb()
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	err := l.BindJSON(&input)
	if err != nil {
		response.Response(l,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Input error please check it!",
			)
		return
	}
	//Whether Username exist or not
	var user model.UserInfos
	db.Where("username = ?", input.Username).First(&user)
	if user.ID == 0 {
		response.Response(l,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"There are no username like this ,please register",
			)
		return
	}
	//Whether password correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		response.Response(l,
			http.StatusBadRequest,
			400,
			nil,
			"Incorrect Password Please check it",
			)
		return
	}
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(l,
			http.StatusInternalServerError,
			500,
			nil,
			"server error",
			)
		return
	}
	tx := db.Begin()
	if tx.Create(&model.Token{UserID: user.ID, Token: token}).RowsAffected != 1 {
		response.Response(l,
			http.StatusInternalServerError,
			500,
			nil,
			"Internet Server Error!",
			)
		tx.Commit()
		return
	}
	response.Success(l,
		gin.H{"token":token},
		"登录成功",
		)
}
func Logout(l *gin.Context) {
	db := common.Getdb()
	token := l.GetString("token")
	db.Model(model.Token{}).Where("token = ?", token).Delete(&model.Token{})
	response.Success(l,
		nil,
		"logout successful",
		)
}
func Info(c *gin.Context) {
	user, _ := c.Get("user")
	response.Success(c,
		gin.H{"user": dto.ToUserDto(user.(model.UserInfos))},
		"Success",
		)
}
