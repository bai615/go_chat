package controllers

import (
	"net/http"
	"fmt"
	"math/rand"
	"go_chat/models"
	"go_chat/services"
	"go_chat/util"
)

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

	loginok := false
	if mobile == "18600000000" && password == "123456" {
		loginok = true
	}

	if loginok {
		// 返回成功
		// data: {"id": 1, "token": "xxx"}
		data := make(map[string]interface{})
		data["id"] = 1
		data["token"] = "test"
		util.RespOk(writer, data, "")

	} else {
		// 返回失败
		util.RespFail(writer, "用户名或者密码错误")
	}
}

var userService services.UserService

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
