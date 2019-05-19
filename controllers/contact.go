package controllers

import (
	"net/http"
	"go_chat/util"
	"go_chat/args"
	"go_chat/services"
)

var contactService services.ContactService

func AddFriend(w http.ResponseWriter, req *http.Request) {
	// 定义一个参数结构体
	/*request.ParseForm()
	mobile := request.PostForm.Get("mobile")
	passwd := request.PostForm.Get("passwd")
	*/
	var arg args.ContactArg
	util.Bind(req, &arg)
	// 调用service
	err := contactService.AddFriend(arg.Userid, arg.Dstid)
	//
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "好友添加成功")
	}
}
