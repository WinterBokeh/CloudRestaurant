package dao

import (
	"CloudRestaurant/model"
	"CloudRestaurant/tool"
	"fmt"
)

type MemberDao struct {
	*tool.Orm
}

func (orm *MemberDao)ValidateSmsCode(phone, code string) *model.SmsCode {
	var sms model.SmsCode

	if _, err := orm.Where(" phone = ? and code = ? ", phone, code).Get(&sms); err != nil {
		fmt.Println(err.Error())
	}
	return &sms
}

func (orm *MemberDao) QueryByPhone(phone string) *model.Member {
	var member model.Member
	if _, err := orm.Where(" mobile = ? ", phone).Get(&member); err != nil {
		fmt.Println(err.Error())
	}
	return &member
}

func (orm *MemberDao) InsertMember(member model.Member) int64 {
	result, err := orm.InsertOne(&member)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func (orm *MemberDao) InsertSme(smsCode *model.SmsCode) int64 {
	result, err := orm.InsertOne(smsCode)
	if err != nil {
		panic(err.Error())
	}
	return result
}
