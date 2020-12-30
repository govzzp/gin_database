package main

import (
	"github.com/gin-gonic/gin"
	"go_free/Controller"
	"go_free/Middleware"
)

func main() {
	r := gin.Default()
	auth := r.Group("/api/auth")
	{
		auth.POST("/register",Controller.Register)
		auth.PUT("/login",Controller.Login)
		auth.GET("/info",Middleware.AuthMiddleware(),Controller.Info)
	}
	panic(r.Run(":8080"))
}
