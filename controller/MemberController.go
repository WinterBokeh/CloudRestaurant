package controller

import (
	"CloudRestaurant/param"
	"CloudRestaurant/service"
	"CloudRestaurant/tool"
	"fmt"
	"github.com/gin-gonic/gin"
)

type MemberController struct { }

func (mc *MemberController) Router(engine *gin.Engine) {
	engine.GET("/api/sendcode", mc.sendSmsCode)

	engine.POST("api/login_sms", mc.smsLogin)

	engine.GET("/api/captcha", mc.captcha)

	engine.POST("/api/vertifycha", mc.vertifycha)
}

func (mc *MemberController) vertifycha(ctx *gin.Context) {
	var captcha tool.CaptchaResult
	err := tool.Decode(ctx.Request.Body, &captcha)
	if err != nil {
		ctx.JSON(200, gin.H{
			"message": "参数解析失败",
		})
		return
	}
	result := tool.VertifyCaptcha(captcha.Id, captcha.VertifyValue)
	if result{
		fmt.Println("验证通过")
	}else {
		fmt.Println("验证失败")
	}
}

func (mc *MemberController) captcha(ctx *gin.Context) {
	// 生成验证码
	tool.GenerateCaptcha(ctx)
}

func (mc *MemberController) smsLogin(ctx *gin.Context) {
	var smsLoginParam param.SmsLoginParam
	err := tool.Decode(ctx.Request.Body, &smsLoginParam)
	if err != nil {
		ctx.JSON(200, gin.H{
			"code": 0,
			"message": err,
		})
	}

	us := service.MemberService{}
	member := us.SmsLogin(smsLoginParam)
	if member != nil {
		ctx.JSON(200, gin.H{
			"code": 1,
			"message": "成功",
			"data": member,
		})
		return
	}
	ctx.JSON(200, gin.H{
		"code": 0,
		"message" : "失败",
	})
}

func (mc *MemberController) sendSmsCode(ctx *gin.Context) {
	phone, exist := ctx.GetQuery("phone")
	if !exist {
		ctx.JSON(200, gin.H{
			"code": 0,
			"message": "参数解析失败",
		})
	}

	ms :=service.MemberService{}
	flag := ms.SendCode(phone)
	if flag {
		ctx.JSON(200, gin.H{
			"code": 1,
			"message" : "发送成功",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 0,
		"message": "发送失败",
	})

}