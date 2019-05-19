package main

import (
	"net/http"
	"io"
)

func main() {
	// 绑定请求和处理函数
	http.HandleFunc("/user/login",
		func(writer http.ResponseWriter,
			request *http.Request) {
				// 数据库操作
				// 逻辑处理
				// restapi json/xml 返回
				io.WriteString(writer,"hello, world")

	})
	// 启动 web 服务器
	http.ListenAndServe(":8080",nil)
}
