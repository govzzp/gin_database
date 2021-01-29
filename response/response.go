package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Response(c *gin.Context,httpStatus int ,code int ,data gin.H,message string)  {
	c.JSON(httpStatus,gin.H{
		"code":code,
		"data":data,
		"message":message,
	})
}
func Success(c *gin.Context,date gin.H,message string)  {
	Response(c,http.StatusOK,200,date,message)
}
func Fail(c *gin.Context,data gin.H ,message string)  {
	Response(c,http.StatusBadRequest,400,data,message)
}
