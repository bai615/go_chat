package controllers

import (
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"net/http"
	"sync"
	"log"
	"fmt"
	"strconv"
)

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
		Conn: conn,
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
		fmt.Printf("recv<=%s", data)
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
