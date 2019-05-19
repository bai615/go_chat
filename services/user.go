package services

import (
	"go_chat/models"
	"errors"
	"fmt"
	"math/rand"
	"go_chat/util"
	"time"
)

type UserService struct {
}

// 注册函数
func (s *UserService) Register(
	mobile,        // 手机号
	plainPassword, // 明文密码
	nickname,      // 昵称
	avatar,        // 头像
	sex string) (user models.User, err error) {

	// 检测手机号是否存在
	userModel := models.User{}
	_, err = DBEngine.Where("mobile=?", mobile).Get(&userModel)
	if nil != err {
		return userModel, err
	}
	// 如果存在则返回提示已经注册
	if userModel.Id > 0 {
		return userModel, errors.New("该手机号已经注册")
	}
	// 否则插入数据库
	userModel.Mobile = mobile
	userModel.Avatar = avatar
	userModel.Nickname = nickname
	userModel.Sex = sex

	userModel.Salt = fmt.Sprintf("%06d", rand.Int31n(10000))
	userModel.Password = util.MakePasswd(plainPassword, userModel.Salt)

	userModel.Createat = time.Now()

	_, err = DBEngine.InsertOne(&userModel)

	// 返回新用户信息
	return userModel, err
}

// 登录函数
func (s *UserService) Login(
	mobile,
	plainPassword string) (user models.User, err error) {

	return user, nil
}
