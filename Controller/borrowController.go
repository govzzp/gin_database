package Controller

import (
	"github.com/gin-gonic/gin"
	"go_free/common"
	"go_free/model"
	"net/http"
)

func Lend(l *gin.Context) {
	db := common.Getdb()
	var input struct {
		Username string `binding:"required" json:"username"`
		ISBN     string `binding:"required" json:"isbn"`
		BookName string `binding:"required" json:"bookname"`
		BookCode string `binding:"required" json:"bookcode"`
	}
	err := l.BindJSON(&input)
	if err != nil {
		l.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Input error please check your input!",
		})
		return
	}
	if len(input.BookCode) != 13 {
		l.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Book_Code must be 13 words!",
		})
		return
	}
	var book model.Book
	db.Table("books").Where("book_code = ?", input.BookCode).First(&book)
	if book.ID == 0 {
		l.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Book is not exist,please check it",
		})
		return
	}
	var user model.UserInfos
	db.Table("user_infos").Where("username = ?", input.Username).First(&user)
	if user.ID == 0 {
		l.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "User not exist please Register user",
		})
		return
	}
	nowBorrow := model.Borrow{
		Username: input.Username,
		ISBN:     input.ISBN,
		BookCode: input.BookCode,
		BookName: input.BookName,
	}
	tx := db.Begin()
	if tx.Create(&nowBorrow).RowsAffected != 1 {
		tx.Rollback()
		l.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Internet Server Error!",
		})
		tx.Commit()
		return
	}
	l.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "borrow book successful!",
	})
}
func Back(b *gin.Context) {
	db := common.Getdb()
	var input struct {
		Username string `binding:"required" json:"username"`
		ISBN     string `binding:"required" json:"isbn"`
		BookName string `binding:"required" json:"bookname"`
		BookCode string `binding:"required" json:"bookcode"`
	}
	err := b.BindJSON(&input)
	if err != nil {
		b.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "input error please check it !",
		})
		return
	}
	if len(input.BookCode) != 13 {
		b.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "BookCode must be 13 words!",
		})
		return
	}
	var book model.Book
	db.Table("books").Where("book_code = ? ", input.BookCode).First(&book)
	if book.ID == 0 {
		b.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Unknow Book Please call the Administrator!",
		})
		return
	}
	var user model.UserInfos
	db.Table("user_infos").Where("username = ?", input.Username).First(&user)
	if user.ID == 0 {
		b.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "The Username is not exist Please Register it !",
		})
		return
	}
	var borrow model.Borrow
	db.Where("book_code = ? ", input.BookCode).First(&borrow)
	if borrow.ID == 0 {
		b.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "You are not borrow this book ,Please call the Administrator!",
		})
		return
	}
	tx := db.Begin()
	if tx.Model(&model.Borrow{}).Where("book_code = ?", input.BookCode).Delete(&model.Borrow{}).RowsAffected != 1 {
		tx.Rollback()
		b.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Database Error!",
		})
		return
	}
	b.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "You Successfully back this book",
		"id":   borrow.ID,
	})
}
