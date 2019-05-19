package main

import (
	"net/http"
	"log"
	"html/template"
	"go_chat/controllers"
)

// 注册模板
func RegisterView() {
	tpl, err := template.ParseGlob("views/**/*")
	if nil != err {
		// 打印错误并直接退出
		log.Fatal(err.Error())
	}

	for _, v := range tpl.Templates() {
		tplName := v.Name()

		http.HandleFunc(tplName, func(writer http.ResponseWriter,
			request *http.Request) {
			tpl.ExecuteTemplate(writer, tplName, nil)
		})
	}
}

func main() {
	// 绑定请求和处理函数
	http.HandleFunc("/user/login", controllers.UserLogin)
	http.HandleFunc("/user/register", controllers.UserRegister)

	// 1、提供静态资源目录支持
	// http.Handle("/", http.FileServer(http.Dir(".")))

	// 2、指定目录的静态文件
	http.Handle("/asset/",
		http.FileServer(http.Dir(".")))

	/*
	// user/login.shtml
	http.HandleFunc("/user/login.shtml", func(writer http.ResponseWriter, request *http.Request) {
		// 解析
		tpl, err := template.ParseFiles("view/user/login.html")
		if nil != err {
			// 打印错误并直接退出
			log.Fatal(err.Error()) // 如果出错直接退出
		}
		tpl.ExecuteTemplate(writer, "/user/login.shtml", nil)
	})

	// user/register.shtml
	http.HandleFunc("/user/register.shtml", func(writer http.ResponseWriter, request *http.Request) {
		// 解析
		tpl, err := template.ParseFiles("view/user/register.html")
		if nil != err {
			// 打印错误并直接退出
			log.Fatal(err.Error()) // 如果出错直接退出
		}
		tpl.ExecuteTemplate(writer, "/user/register.shtml", nil)
	})
	*/

	RegisterView()

	// 启动 web 服务器
	http.ListenAndServe(":8080", nil)
}
