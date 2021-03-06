package controllers

import (
	"net/http"
	"go_chat/util"
	"go_chat/args"
	"go_chat/services"
	"go_chat/models"
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

func LoadFriend(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	// 如果这个用的上,那么可以直接
	util.Bind(req, &arg)

	users := contactService.SearchFriend(arg.Userid)
	util.RespOkList(w, users, len(users))
}

func LoadCommunity(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	// 如果这个用的上,那么可以直接
	util.Bind(req, &arg)
	comunitys := contactService.SearchCommunity(arg.Userid)
	util.RespOkList(w, comunitys, len(comunitys))
}

func CreateCommunity(w http.ResponseWriter, req *http.Request) {
	var arg models.Community
	// 如果这个用的上,那么可以直接
	util.Bind(req, &arg)
	com, err := contactService.CreateCommunity(arg)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, com, "")
	}
}

func JoinCommunity(w http.ResponseWriter, req *http.Request) {
	var arg args.ContactArg
	println(arg.Userid)
	// 如果这个用的上,那么可以直接
	util.Bind(req, &arg)
	println("userid: " , arg.Userid)
	println("dstid: ", arg.Dstid)
	err := contactService.JoinCommunity(arg.Userid, arg.Dstid)
	// todo 刷新用户的群组信息
	AddGroupId(arg.Userid, arg.Dstid)
	if err != nil {
		util.RespFail(w, err.Error())
	} else {
		util.RespOk(w, nil, "")
	}
}
