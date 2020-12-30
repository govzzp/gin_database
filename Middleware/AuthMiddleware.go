package Middleware

import (
	"github.com/gin-gonic/gin"
	"go_free/common"
	"go_free/model"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Insufficient permissions",
			})
			c.Abort()
			return
		}
		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Insufficient permissions",
			})
			c.Abort()
			return
		}
		UserID := claims.UserID
		db := common.Getdb()
		var user model.UserInfos
		db.First(&user, UserID)
		if UserID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "Insufficient permissions",
			})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}
