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
		ISBN       string `binding:"required" json:"isbn"`
		BookName   string `binding:"required" json:"bookname"`
		Author     string `binding:"required" json:"author"`
		Press      string `binding:"required" json:"press"`
		Category   string `binding:"required" json:"category"`   //种类
		GetBooknum string `binding:"required" json:"getbooknum"`
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
	newBook :=model.Book{
		ISBN:     input.ISBN,
		BookName: input.BookName,
		Author: input.Author,
		Press: input.Press,
		Category: input.Category,
		GetBooknum : input.GetBooknum,
	}

	tx := db.Begin()
	if  tx.Create(&newBook).RowsAffected != 1{
		add.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "Internet Server Error!",
		})
		tx.Commit()
		return
	}
	add.JSON(http.StatusCreated,gin.H{
		"code":201,
		"msg":"Add book Successful",
	})
}
func Changebook() {

}
func Deletebook() {

}
func Allbooks() {

}
