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

func (ms *MemberService) UploadAvatar(userId int64, fileName string) string {
	memberDao := dao.MemberDao{tool.DbEngine}
	result := memberDao.UpdateMemberAvatar(userId, fileName)
	if result == 0 {
		return ""
	}

	return fileName
}

func (ms *MemberService) PwdLogin(nameAndPwd param.NameAndPassword) *model.Member {
	md := dao.MemberDao{tool.DbEngine}
	//检查是否已存在
	member := md.Query(nameAndPwd.UserName)
	if member.Id != 0 {
		if member.Password != nameAndPwd.Password {
			return new(model.Member)
		}
		return member
	}

	user := model.Member{}
	user.UserName = nameAndPwd.UserName
	user.Password = nameAndPwd.Password
	user.RegisterTime = time.Now().Unix()

	user.Id = md.InsertMember(user)
	return &user
}

//业务逻辑
func (ms *MemberService) SmsLogin(loginparam param.SmsLoginParam) *model.Member {
	md := dao.MemberDao{tool.DbEngine}
//	fmt.Println(loginparam.Phone, "hahahahaha")
	//从sms_code查看给出的code和phone是否匹配
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
