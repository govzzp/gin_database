package Controller

import (
	"github.com/gin-gonic/gin"
	"go_free/common"
	"go_free/model"
	"go_free/response"
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
		response.Response(l,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Input error please check your input!",
			)
		return
	}
	if len(input.BookCode) != 13 {
		response.Response(l,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Book_Code must be 13 words!",
			)
		return
	}
	var book model.Book
	db.Table("books").Where("book_code = ?", input.BookCode).First(&book)
	if book.ID == 0 {
		response.Response(l,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Book is not exist,please check it",
			)
		return
	}
	if input.BookName != book.BookName {
		response.Response(l,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Error BookName Please Check it!",
		)
		return
	}
	var user model.UserInfos
	db.Table("user_infos").Where("username = ?", input.Username).First(&user)
	if user.ID == 0 {
		response.Response(l,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"User not exist please Register user",
		)
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
		nil,
		"borrow book successful!",
		)
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
		response.Response(b,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"input error please check it !",
		)

		return
	}
	if len(input.BookCode) != 13 {
		response.Response(b,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"BookCode must be 13 words!",
		)
		return
	}
	var book model.Book
	db.Table("books").Where("book_code = ? ", input.BookCode).First(&book)
	if book.ID == 0 {
		response.Response(b,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Unknow Book Please call the Administrator!",
		)
		return
	}
	if input.BookName != book.BookName {
		response.Response(b,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"Unknow Book Please call the Administrator!",
		)
		return
	}
	var user model.UserInfos
	db.Table("user_infos").Where("username = ?", input.Username).First(&user)
	if user.ID == 0 {
		response.Response(b,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"The Username is not exist Please Register it !",
		)
		return
	}
	var borrow model.Borrow
	db.Where("book_code = ? ", input.BookCode).First(&borrow)
	if borrow.ID == 0 {
		response.Response(b,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"You are not borrow this book ,Please call the Administrator!",
		)
		return
	}
	tx := db.Begin()
	if tx.Model(&model.Borrow{}).Where("book_code = ?", input.BookCode).Delete(&model.Borrow{}).RowsAffected != 1 {
		tx.Rollback()
		response.Response(b,
			http.StatusInternalServerError,
			500,
			nil,
			"Database Error!",
		)
		return
	}
	response.Success(b,
		gin.H{"id":borrow.ID},
		"You Successfully back this book",
		)
}
