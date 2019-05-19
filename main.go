package main

import (
	"net/http"
)

func main() {
	// 绑定请求和处理函数
	http.HandleFunc("/user/login",
		func(writer http.ResponseWriter,
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

			var str string
			if loginok {
				// 返回成功 JSON
				str = `{"code": 0, "data": {"id": 1, "token": "test"}}`

			} else {
				// 返回失败 JSON
				str = `{"code": -1, "message": "用户名或者密码错误"}`
			}

			// 设置 header 为 JSON，默认为 text/html，所以特别指出返回为 application/json
			writer.Header().Set("Content-Type", "application/json")
			// 设置 200 状态
			writer.WriteHeader(http.StatusOK)
			// 输出
			writer.Write([]byte(str))

			// 如何返回 JSON
			//io.WriteString(writer, "hello, world")

		})
	// 启动 web 服务器
	http.ListenAndServe(":8080", nil)
}
