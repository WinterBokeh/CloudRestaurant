package controller

import "github.com/gin-gonic/gin"

type HelloController struct { }

func (hello *HelloController) Router(engine *gin.Engine) {
	engine.GET("/hello", hello.Hello)
}

func (hello *HelloController) Hello(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello",
		"code": "ok",
	})
}