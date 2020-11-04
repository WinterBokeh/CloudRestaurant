package service

import (
	"CloudRestaurant/dao"
	"CloudRestaurant/model"
	"CloudRestaurant/param"
	"CloudRestaurant/tool"
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/gin-gonic/gin"
	"math/rand"
	"time"
)

type MemberService struct { }


//业务逻辑
func (ms *MemberService) SmsLogin(loginparam param.SmsLoginParam) *model.Member {
	md := dao.MemberDao{tool.DbEngine}
//	fmt.Println(loginparam.Phone, "hahahahaha")
	sms := md.ValidateSmsCode(loginparam.Phone, loginparam.Code)
	if sms.Id == 0 {
		return nil
	}

	member := md.QueryByPhone(loginparam.Phone)
	if member.Id != 0 {
		return member
	}

	user := model.Member{}
	user.UserName = loginparam.Phone
	user.Mobile = loginparam.Phone
	user.RegisterTime = time.Now().Unix()

	user.Id = md.InsertMember(user)

	return &user
}

func (ms *MemberService) SendCode(phone string) bool {
	//产生一个验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	//调用阿里云sdk
	config := tool.GetConfig().Sms
	client, err := dysmsapi.NewClientWithAccessKey(config.RegionId, config.AppKey, config.AppSecret)
	if err != nil {
		return false
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = config.SignName
	request.TemplateCode = config.TemplateCode
	request.PhoneNumbers = phone

	par, err := json.Marshal(gin.H{
		"code": code,
	})

	request.TemplateParam = string(par)

	response, err := client.SendSms(request)
	fmt.Println(response)

	if err != nil {
		return false
	}

	if response.Code == "OK" {
		smsCode := model.SmsCode{Phone: phone, Code: code, BizId: response.BizId, CreateTime: time.Now().Unix()}
		Dao := dao.MemberDao{tool.DbEngine}
		result := Dao.InsertSme(&smsCode)
		return  result > 0
	}
	//接收返回结果，并判断状态
	return false
}
