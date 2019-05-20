package controllers

import (
	"net/http"
	"os"
	"go_chat/util"
	"strings"
	"fmt"
	"time"
	"math/rand"
	"io"
)

func init() {
	os.MkdirAll("./mnt", os.ModePerm) // 0777
}

func Upload(w http.ResponseWriter, r *http.Request) {
	UploadLocal(w, r)
}

// 1、存储位置 ./mnt，需要确保已经创建好
// 2、url 格式 /mnt/xxxx.png 需要确保网络能访问 /mnt/
func UploadLocal(writer http.ResponseWriter, request *http.Request) {
	// todo 获得上传的源文件s
	srcFile, head, err := request.FormFile("file")
	if nil != err {
		util.RespFail(writer, err.Error())
		return
	}
	// todo 创建一个新文件d
	suffix := ".png"
	// 如果前端文件名称包含后缀 xxx.png
	oFileName := head.Filename
	tmp := strings.Split(oFileName, ".")
	if len(tmp) > 1 {
		suffix = "." + tmp[len(tmp)-1]
	}
	// 如果前端指定 fileType
	fileType := request.FormValue("fileType")
	if len(fileType) > 0 {
		suffix = fileType
	}
	fileName := fmt.Sprintf("%d%04d%s",
		time.Now().Unix(),
		rand.Int31(),
		suffix)
	// 创建新文件
	dstFile, err := os.Create("./mnt/" + fileName)
	if nil != err {
		util.RespFail(writer, err.Error())
		return
	}
	// todo 将源文件内容 copy 到新文件
	_, err = io.Copy(dstFile, srcFile)
	if nil != err {
		util.RespFail(writer, err.Error())
		return
	}

	// todo 将新文件路径转换成 URL 地址

	url := "/mnt/" + fileName

	// todo 响应到前端
	util.RespOk(writer, url, "")
}
