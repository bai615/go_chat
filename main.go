package main

import (
	"net/http"
	"encoding/json"
	"log"
)

func userLogin(writer http.ResponseWriter,
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
		Resp(writer, 0, data, "")

	} else {
		// 返回失败
		Resp(writer, -1, nil, "用户名或者密码错误")
	}
}

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message""`
	Data    interface{} `json:"data,omitempty"` // omitempty：data 为 null 时，不显示
}

func Resp(writer http.ResponseWriter, code int, data interface{}, message string) {
	// 设置 header 为 JSON，默认为 text/html，所以特别指出返回为 application/json
	writer.Header().Set("Content-Type", "application/json")
	// 设置 200 状态
	writer.WriteHeader(http.StatusOK)

	// 定义一个结构体
	responseData := ResponseData{
		Code:    code,
		Message: message,
		Data:    data,
	}
	// 将结构体转化为 JSON 字符串
	ret, err := json.Marshal(responseData)
	if err != nil {
		log.Println(err.Error())
	}
	// 输出
	writer.Write(ret)
}

func main() {
	// 绑定请求和处理函数
	http.HandleFunc("/user/login", userLogin)
	// 启动 web 服务器
	http.ListenAndServe(":8080", nil)
}
