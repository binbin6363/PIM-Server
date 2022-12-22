package client

import (
	"PIM_Server/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"runtime/debug"
	"time"
)

// Login 用户登录
type Login struct {
	PlatformID PlatformType
	Uid        int64
	Client     *Client
}

// GetKey 获取 key
func (l *Login) GetKey() (key string) {
	key = GetUserKey(l.PlatformID, l.Uid)

	return
}

type PlatformType int

const (
	PlatformPc PlatformType = iota + 1
	PlatformWeb
	PlatformIos
	PlatformAndroid
)

type Client struct {
	Socket        *websocket.Conn // websocket的连接对象
	Addr          string          // 对端的地址
	Send          chan []byte     // 待发送的数据
	PlatformId    PlatformType    // 登录的平台Id 1 pc app/2 web/3 ios
	Uid           int64           // 用户Id，用户登录以后才有
	ConnTime      int64           // 首次连接时间
	HeartbeatTime int64           // 用户上次心跳时间
	LoginTime     int64           // 登录时间 登录以后才有
}

func New(addr string, socket *websocket.Conn) (client *Client) {
	client = &Client{
		Addr:     addr,
		Socket:   socket,
		Send:     make(chan []byte, 100),
		ConnTime: time.Now().Unix(),
	}
	client.HeartbeatTime = client.ConnTime

	return
}

// Close 客户端连接关闭
func (c *Client) Close() {
	close(c.Send)
	if err := c.Socket.Close(); err != nil {
		fmt.Println("conn close failed,", c.Addr, c.PlatformId, c.Uid, c.ConnTime, c.HeartbeatTime, c.LoginTime)
	} else {
		fmt.Println("conn close success,", c.Addr, c.PlatformId, c.Uid, c.ConnTime, c.HeartbeatTime, c.LoginTime)
	}

}

// SendRawMsg 发送数据
func (c *Client) SendRawMsg(msg []byte) {

	if c == nil {
		log.Infof("client nil")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendRawMsg stop:", r, string(debug.Stack()))
		}
	}()

	c.Send <- msg
}

// SendJsonMsg 发送json数据
func (c *Client) SendJsonMsg(obj interface{}) error {

	if c == nil {
		log.Infof("[ERROR] client nil")
		return errors.New("client not exist")
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendJsonMsg stop:", r, string(debug.Stack()))
		}
	}()

	if data, err := json.Marshal(obj); err != nil {
		log.Infof("[ERROR] send error, json Marshal failed, err:%+v", err)
		return err
	} else {
		log.Infof("send %s", string(data))
		c.Send <- data
		return nil
	}
}

func (c *Client) Run() {
	for {
		select {
		case data := <-c.Send:
			// 发送数据
			log.Infof("real send data to socket")
			c.Socket.WriteMessage(websocket.TextMessage, data)
		}
	}
}
