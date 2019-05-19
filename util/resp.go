package util

import (
	"net/http"
	"encoding/json"
	"log"
)

type ResponseData struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"` // omitempty：data 为 null 时，不显示
	Rows    interface{} `json:"rows,omitempty"`
	Total   interface{} `json:"total,omitempty"`
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

func RespFail(writer http.ResponseWriter, message string) {
	Resp(writer, -1, nil, message)
}

func RespOk(writer http.ResponseWriter, data interface{}, message string) {
	Resp(writer, 0, data, message)
}

func RespList(w http.ResponseWriter, code int, data interface{}, total interface{}) {

	w.Header().Set("Content-Type", "application/json")
	// 设置200状态
	w.WriteHeader(http.StatusOK)
	// 输出
	// 定义一个结构体
	// 满足某一条件的全部记录数目
	// 测试 100
	// 20
	responseData := ResponseData{
		Code:  code,
		Rows:  data,
		Total: total,
	}
	// 将结构体转化成JSOn字符串
	ret, err := json.Marshal(responseData)
	if err != nil {
		log.Println(err.Error())
	}
	// 输出
	w.Write(ret)
}

func RespOkList(w http.ResponseWriter, lists interface{}, total interface{}) {
	// 分页数目,
	RespList(w, 0, lists, total)
}
