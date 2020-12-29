package main

import (
	"github.com/gin-gonic/gin"
	"go_free/Controller"
)




func main() {
	r := gin.Default()
	r.POST("/user", Controller.Register)
	r.PUT("/user", Controller.Login)
	panic(r.Run(":8080"))
}
