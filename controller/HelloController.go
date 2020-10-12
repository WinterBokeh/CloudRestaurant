package controller

import (
	"github.com/gin-gonic/gin"
)

type HelloController struct {
}

func (hello *HelloController) Router(engin *gin.Engine) {
	engin.GET("/hello", hello.Hello)
}

func (hello *HelloController) Hello(context *gin.Context) {
	context.JSON(200, gin.H{
		"message": "ok",
	})
}
