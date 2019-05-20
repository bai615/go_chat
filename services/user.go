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

	userModel.Token = fmt.Sprintf("%08d", rand.Int31())

	userModel.Createat = time.Now()

	_, err = DBEngine.InsertOne(&userModel)

	// 返回新用户信息
	return userModel, err
}

// 登录函数
func (s *UserService) Login(
	mobile,
	plainPassword string) (user models.User, err error) {

	// 通过手机号查询用户
	userModel := models.User{}
	_, err = DBEngine.Where("mobile=?", mobile).Get(&userModel)
	if nil != err {
		return userModel, err
	}

	if userModel.Id == 0 {
		return userModel, errors.New("该用户不存在")
	}

	// 对比密码
	loginok := util.ValidatePasswd(plainPassword, userModel.Salt, userModel.Password)
	if !loginok {
		return userModel, errors.New("用户名或者密码错误")
	}

	// 刷新 token，安全
	str := fmt.Sprintf("%d", time.Now().Unix())
	token := util.MD5Encode(str)
	userModel.Token = token

	DBEngine.ID(userModel.Id).Cols("token").Update(&userModel)

	return userModel, nil
}

// 查找某个用户
func (s *UserService) Find(
	userId int64) (user models.User) {

	// 首先通过手机号查询用户
	tmp := models.User{

	}
	DBEngine.ID(userId).Get(&tmp)
	return tmp
}
