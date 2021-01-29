package Controller

import (
	"github.com/gin-gonic/gin"
	"go_free/common"
	"go_free/model"
	"go_free/response"
	"net/http"
)

func Addbook(add *gin.Context) {
	db := common.Getdb()
	var input struct {
		BookCode   string `binding:"required" json:"bookcode"`
		ISBN       string `binding:"required" json:"isbn"`
		BookName   string `binding:"required" json:"bookname"`
		Author     string `binding:"required" json:"author"`
		Press      string `binding:"required" json:"press"`    //出版社
		Category   string `binding:"required" json:"category"` //种类
		GetBooknum string `binding:"required" json:"getbooknum"`
		Position   string `binding:"required" json:"position"`
	}
	err := add.BindJSON(&input)
	if err != nil {
		response.Response(add,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"input book info error please check your input",
			)
		return
	}
	if len(input.ISBN) != 13 {
		response.Response(add,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"ISBN number must be 13 words",
		)
		return
	}
	var book model.Book
	db.Where("book_code = ?", input.BookCode).First(&book)
	if book.ID != 0 {
		response.Response(add,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"This Code of book is already be used ,please use another",
		)
		return
	}
	newBook := model.Book{
		BookCode:   input.BookCode,
		ISBN:       input.ISBN,
		BookName:   input.BookName,
		Author:     input.Author,
		Press:      input.Press,
		Category:   input.Category,
		GetBooknum: input.GetBooknum,
		Position:   input.Position,
	}

	tx := db.Begin()
	if tx.Create(&newBook).RowsAffected != 1 {
		response.Response(add,
			http.StatusInternalServerError,
			500,
			nil,
			"Internet Server Error!",
		)
		tx.Commit()
		return
	}
	response.Success(add,nil,"Add book Successful")
}
func Changebook(c *gin.Context) {
	db := common.Getdb()
	var input struct {
		BookCode   string `binding:"required" json:"bookcode"`
		ISBN       string `binding:"required" json:"isbn"`
		BookName   string `binding:"required" json:"bookname"`
		Author     string `binding:"required" json:"author"`
		Press      string `binding:"required" json:"press"`    //出版社
		Category   string `binding:"required" json:"category"` //种类
		GetBooknum string `binding:"required" json:"getbooknum"`
		Position   string `binding:"required" json:"position"`
	}
	err := c.BindJSON(&input)
	if err != nil {
		response.Response(c,
			http.StatusUnprocessableEntity,
			422,
			nil,
			"input error please check your input !",
		)
		return
	}
	updateBook := model.Book{
		BookCode:   input.BookCode,
		ISBN:       input.ISBN,
		BookName:   input.BookName,
		Author:     input.Author,
		Press:      input.Press,
		Category:   input.Category,
		GetBooknum: input.GetBooknum,
		Position:   input.Position,
	}
	tx := db.Begin()
	if tx.Where("book_code = ?", input.BookCode).UpdateColumns(&updateBook).RowsAffected != 1 {
		db.Rollback()
		response.Response(c,
			http.StatusInternalServerError,
			500,
			nil,
			"Internet Server Error!",
		)
		return
	}
	tx.Commit()
	var book model.Book
	db.Where("book_code = ?", input.BookCode).First(&book)
	response.Success(c,gin.H{"id":book.ID},"Book Update OK ")
}
func Deletebook(d *gin.Context) {
	db := common.Getdb()
	var input struct {
		BookCode string `binding:"required" json:"bookcode"`
	}
	err := d.BindJSON(&input)
	if err != nil {
		response.Response(d,
			http.StatusUnprocessableEntity,
			422,nil,
			"Input Error Please Check it!",
			)
		return
	}

	tx := db.Begin()
	if tx.Model(&model.Book{}).Where("book_code = ?", input.BookCode).Delete(&model.Book{}).RowsAffected != 1 {
		tx.Rollback()
		response.Response(d,
			http.StatusInternalServerError,
			500,
			nil,
			"Database Error!",
			)
		return
	}
	response.Success(d,
		nil,
		"Delete book Successful!",
		)
}

func Allbooks(a *gin.Context) {
	var books []model.Book
	db := common.Getdb()
	db.Model(&model.Book{}).Find(&books)
	response.Success(a,
		gin.H{"data": books},
		"Success",
		)
}
