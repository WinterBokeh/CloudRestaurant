package tool

import "github.com/gin-gonic/gin"

func Success(ctx *gin.Context, v interface{})  {
	ctx.JSON(200, gin.H{
		"code": 0,
		"smg": "成功",
		"data": v,
	})
}

func Faild(ctx *gin.Context, v interface{})  {
	ctx.JSON(200, gin.H{
		"code": 1,
		"smg": "失败",
		"data": v,
	})
}
