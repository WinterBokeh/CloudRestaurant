package controller

import (
	"CloudRestaurant/model"
	"CloudRestaurant/param"
	"CloudRestaurant/service"
	"CloudRestaurant/tool"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

type MemberController struct { }

func (mc *MemberController) Router(engine *gin.Engine) {
	engine.Static("/uploadfile", "./uploadfile")

	engine.GET("/api/sendcode", mc.sendSmsCode)

	engine.POST("api/login_sms", mc.smsLogin)

	engine.GET("/api/captcha", mc.captcha)

//	engine.POST("/api/vertifycha", mc.vertifycha)
	engine.POST("/api/login_pwd", mc.login)

	engine.POST("/api/upload/avator", mc.uploadAvator)
}

func (mc *MemberController) uploadAvator(ctx *gin.Context) {
	//解析上传的文件 file， user_id
	userId := ctx.PostForm("user_id")
//	fmt.Println("阿巴阿巴", userId)
	file, err := ctx.FormFile("avatar")
	if err != nil {
		tool.Faild(ctx, "参数解析失败")
		return
	}

	//判断用户是否已经登录
	sess := tool.GetSession(ctx, "user_" + userId)
//	fmt.Println("XXXXXXXXXXXX", "user_" + userId, "XXXXXXXXXXX")
	if sess == nil {
		tool.Faild(ctx, "参数不合法")
		return
	}

	var member model.Member
	json.Unmarshal(sess.([]byte), &member)

	//把file保存到本地
	fileName := "./uploadfile/" + strconv.FormatInt( time.Now().Unix() , 10 ) + file.Filename
	err = ctx.SaveUploadedFile(file, fileName)
	if err != nil {
		fmt.Println(err)
		tool.Faild(ctx, "头像更新失败")
		return
	}

	//把保存后的本地路径保存到路径表中头像字段
	ms := service.MemberService{}
	path :=  ms.UploadAvatar(member.Id, fileName[1:])
	if path != "" {
		tool.Success(ctx, "http://localhost:8090" + path)
		return
	}

	tool.Faild(ctx, "上传失败")
	//返回结果
}

func (mc *MemberController) login(ctx *gin.Context) {
	//获取表单，解析表单
	var loginByNP  param.NameAndPassword
	err := tool.Decode(ctx.Request.Body, &loginByNP)
	if err != nil {
		tool.Faild(ctx, "参数解析失败")
		return
	}
//	fmt.Println(loginByNP.UserName, "阿巴阿巴")

	//验证验证码

	flag := tool.VertifyCaptcha(loginByNP.Id, loginByNP.Code)
	if flag == false {
		tool.Faild(ctx, "验证码错误")
		return
	}

	//插入数据库

	us := service.MemberService{}
	member := us.PwdLogin(loginByNP)

	if member.Id != 0 {
		//保存session
		sess, _ := json.Marshal(member)
//		fmt.Println("test1", member.Id)
//		fmt.Println("test2", strconv.FormatInt(member.Id, 10))
//		fmt.Println("XXXXXXXXXXXX", "user_" + strconv.FormatInt(member.Id, 10), "XXXXXXXXXXX")
		err := tool.SetSession(ctx, "user_" + strconv.FormatInt(member.Id, 10), sess)
		if err != nil {
			tool.Faild(ctx, "登录失败")
			return
		}
		tool.Success(ctx, &member)
		return
	}
	tool.Faild(ctx, "密码错误")
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
		sess, _ := json.Marshal(member)
		err := tool.SetSession(ctx, "user_" + strconv.FormatInt(member.Id, 10), sess)
		if err != nil {
			tool.Faild(ctx, "登陆失败")
		}
		tool.Success(ctx, member)
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