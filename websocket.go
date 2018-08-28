
eckage main

import (
	"net/http"

	"sync"

	"fmt"
	"time"

	"errors"

	"log"

	"github.com/gorilla/websocket"
)

//http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	//允许所有cors跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

//客户端读写消息
type wsMessage struct {
	messageType int
	data        []byte
}

//客户端链接
type wsConnection struct {
	wsSocket *websocket.Conn //底层的websocket
	inChan   chan *wsMessage //读队列
	outChan  chan *wsMessage //写队列

	mutex     sync.Mutex //避免重复关闭管道
	isClosed  bool
	closeChan chan byte //关闭通知
}

func main() {
	port := 7777
	http.HandleFunc("/ws", wsHandler)
	log.Printf("websocket server on :%d", port)
	http.ListenAndServe(fmt.Sprintf("0.0.0.0.:%d", port), nil)
}

func wsHandler(resp http.ResponseWriter, req *http.Request) {
	//应答客户端告知升级链接为websocket
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}

	wsConn := &wsConnection{
		wsSocket:  wsSocket,
		inChan:    make(chan *wsMessage, 1000),
		outChan:   make(chan *wsMessage, 1000),
		closeChan: make(chan byte),
		isClosed:  false,
	}

	log.Println("RemoteAddr: ", wsConn.wsSocket.RemoteAddr())
	log.Println("LocalAddr: ", wsConn.wsSocket.LocalAddr())

	//处理器
	go wsConn.proLoop()
	//读协程
	go wsConn.wsReadLoop()
	//写协程
	go wsConn.wsWriteLoop()
}

func (wsConn *wsConnection) proLoop() {
	//启动一个gouroutine发送心跳
	go func() {
		for {
			time.Sleep(2 * time.Second)
			if err := wsConn.wsWrite(websocket.TextMessage, []byte("heartbeat from server")); err != nil {
				fmt.Println("heartbeat fail")
				wsConn.wsClose()
				break
			}
		}
	}()

	// 这是一个同步处理模型，如果希望并行处理可以每个请求一个gorutine, 注意控制并发goroutine的数量
	for {
		msg, err := wsConn.wsRead()
		if err != nil {
			fmt.Println("read fail")
			break
		}
		fmt.Println(string(msg.data))
		if err = wsConn.wsWrite(msg.messageType, msg.data); err != nil {
			fmt.Println("write fail")
			break
		}
	}
}
func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
	select {
	case wsConn.outChan <- &wsMessage{messageType, data}:
	case <-wsConn.closeChan:
		return errors.New("websocket closed")
	}
	return nil
}
func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
	select {
	case msg := <-wsConn.inChan:
		return msg, nil
	case <-wsConn.closeChan:
	}
	return nil, errors.New("websocket closed")
}
func (wsConn *wsConnection) wsClose() {
	wsConn.wsSocket.Close()

	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		wsConn.isClosed = true
		close(wsConn.closeChan)
	}
}
func (wsConn *wsConnection) wsReadLoop() {
	for {
		//读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			wsConn.wsClose()
			break
		}
		req := &wsMessage{
			msgType,
			data,
		}
		//放入队列
		select {
		case wsConn.inChan <- req:
		case <-wsConn.closeChan:
			break
		}
	}
}
func (wsConn *wsConnection) wsWriteLoop() {
	for {
		select {
		//取一个应答
		case msg := <-wsConn.outChan:
			if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
				goto error
			}
		case <-wsConn.closeChan:
			goto closed
		}
	}
error:
	wsConn.wsClose()
closed:
}
