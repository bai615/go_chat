package controllers

import (
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"net/http"
	"sync"
	"log"
	"fmt"
	"strconv"
	"encoding/json"
)

const (
	CMD_SINGLE_MSG = 10
	CMD_ROOM_MSG   = 11
	CMD_HEART      = 0
)

type Message struct {
	Id      int64  `json:"id,omitempty" form:"id"`           // 消息ID
	Userid  int64  `json:"userid,omitempty" form:"userid"`   // 谁发的
	Cmd     int    `json:"cmd,omitempty" form:"cmd"`         // 群聊还是私聊
	Dstid   int64  `json:"dstid,omitempty" form:"dstid"`     // 对端用户ID/群ID
	Media   int    `json:"media,omitempty" form:"media"`     // 消息按照什么样式展示
	Content string `json:"content,omitempty" form:"content"` // 消息的内容
	Pic     string `json:"pic,omitempty" form:"pic"`         // 预览图片
	Url     string `json:"url,omitempty" form:"url"`         // 服务的URL
	Memo    string `json:"memo,omitempty" form:"memo"`       // 简单描述
	Amount  int    `json:"amount,omitempty" form:"amount"`   // 其他和数字相关的
}

type Node struct {
	Conn *websocket.Conn
	// 并行转串行, Conn 为 IO 型资源，有忙和闲的概念，存在竞争关系。一个 Conn 正在写数据时，另一个 Conn 接入数据会混乱。并行转串行后，数据为顺序型。
	DataQueue chan []byte
	GroupSets set.Interface
}

// 映射关系表
var clientMap map[int64]*Node = make(map[int64]*Node, 0)

// 读写锁
var rwlocker sync.RWMutex

// ws://127.0.0.1:8080/chat?id=1&token=xxxx
func Chat(writer http.ResponseWriter, request *http.Request) {

	// todo 检验接入是否合法
	// checkToken(userId int64, token string)
	query := request.URL.Query() // 获取 URL 中参数
	id := query.Get("id")
	token := query.Get("token")
	// 将字符串转为 int64 型
	userId, _ := strconv.ParseInt(id, 10, 64)
	isvalida := checkToken(userId, token)
	// 如果 isvalida = true 继续
	// isvalida = false 返回

	// 自动升级
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isvalida
		},
	}).Upgrade(writer, request, nil)
	if nil != err {
		log.Println(err.Error())
		return
	}

	// todo 获得 conn
	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	// todo userid 和 node 形成绑定关系
	rwlocker.Lock()
	clientMap[userId] = node
	rwlocker.Unlock()

	// todo 完成发送逻辑，con
	go sendproc(node)
	// todo 完成接收逻辑
	go recvproc(node)

	// 发送信息
	sendMsg(userId, []byte("hello, world!"))
}

// 发送协程
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// 接收协程
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}
		// todo 对data进一步处理
		dispatch(data)
		fmt.Printf("recv<=%s", data)
	}
}

// 后端调度逻辑处理
func dispatch(data []byte) {
	// todo 解析 data 为 message
	msg := Message{}
	err := json.Unmarshal(data, &msg)
	if nil != err {
		log.Println(err.Error())
		return
	}
	// todo 根据 cmd 对逻辑进行处理
	switch msg.Cmd {
	case CMD_SINGLE_MSG: // 单聊
		sendMsg(msg.Dstid, data)
	case CMD_ROOM_MSG: // 群聊
		// todo 群聊转发逻辑
	case CMD_HEART: // 心跳
		// todo 一般啥也不做
	}
}

// todo 发送消息
func sendMsg(userId int64, msg []byte) {
	rwlocker.RLock()
	node, ok := clientMap[userId]
	rwlocker.RUnlock()
	if ok {
		node.DataQueue <- msg
	}
}

// 检测是否有效
func checkToken(userId int64, token string) bool {
	// 从数据库里面查询并比对
	user := userService.Find(userId)
	return user.Token == token
}
