package controllers

import (
	"net/http"
	"fmt"
	"math/rand"
	"go_chat/models"
	"go_chat/services"
	"go_chat/util"
)

var userService services.UserService

func UserLogin(writer http.ResponseWriter,
	request *http.Request) {
	// 数据库操作
	// 逻辑处理
	// restapi json/xml 返回
	// 1.获取前端传递的参数
	// mobile, password
	// 解析参数
	// 如何获得参数
	// 解析参数

	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	password := request.PostForm.Get("password")

	user, err := userService.Login(mobile, password)
	if nil != err {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}
}

func UserRegister(writer http.ResponseWriter,
	request *http.Request) {

	request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	plainPassword := request.PostForm.Get("password")
	nickname := fmt.Sprintf("user%06d", rand.Int31())
	avatar := ""
	sex := models.SEX_UNKNOW

	user, err := userService.Register(mobile, plainPassword, nickname, avatar, sex)
	if nil != err {
		util.RespFail(writer, err.Error())
	} else {
		util.RespOk(writer, user, "")
	}
}
