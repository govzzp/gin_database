package Controller

import (
	"github.com/gin-gonic/gin"
	"go_free/common"
	"go_free/model"
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
		add.JSON(http.StatusUnprocessableEntity, gin.H{
			"code":    422,
			"message": "input book info error please check your input",
		})
		return
	}
	if len(input.ISBN) != 13 {
		add.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "ISBN number must be 13 words",
		})
		return
	}
	var book model.Book
	db.Where("book_code = ?", input.BookCode).First(&book)
	if book.ID != 0 {
		add.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "This Code of book is already be used ,please use another",
		})
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
		add.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Internet Server Error!",
		})
		tx.Commit()
		return
	}
	add.JSON(http.StatusCreated, gin.H{
		"code": 201,
		"msg":  "Add book Successful",
	})
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
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "input error please check your input !",
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Internet Server Error!",
		})
		return
	}
	tx.Commit()
	var book model.Book
	db.Where("book_code = ?", input.BookCode).First(&book)
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Book Update OK ",
		"id":   book.ID,
	})
}
func Deletebook(d *gin.Context) {
	db := common.Getdb()
	var input struct {
		BookCode string `binding:"required" json:"bookcode"`
	}
	err := d.BindJSON(&input)
	if err != nil {
		d.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "Input Error Please Check it!",
		})
		return
	}

	tx := db.Begin()
	if tx.Model(&model.Book{}).Where("book_code = ?", input.BookCode).Delete(&model.Book{}).RowsAffected != 1 {
		tx.Rollback()
		d.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Database Error!",
		})
		return
	}
	d.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "Delete book Successful!",
	})
}

//func Allbooks(a *gin.Context) {
//
//	db := common.Getdb()
//	result := db.Find(&model.Book{},n)
//	a.JSON(http.StatusOK,gin.H{
//		"code":200,
//
//	})
//}
