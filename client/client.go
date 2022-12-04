package client

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"runtime/debug"
	"time"
)

// Login 用户登录
type Login struct {
	PlatformID uint32
	Uid        int64
	Client     *Client
}

// GetKey 获取 key
func (l *Login) GetKey() (key string) {
	key = GetUserKey(l.PlatformID, l.Uid)

	return
}

type Client struct {
	Socket        *websocket.Conn // websocket的连接对象
	Addr          string          // 对端的地址
	Send          chan []byte     // 待发送的数据
	PlatformId    uint32          // 登录的平台Id 1 pc app/2 web/3 ios
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

// SendMsg 发送数据
func (c *Client) SendMsg(msg []byte) {

	if c == nil {
		log.Printf("client nil")
		return
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("SendMsg stop:", r, string(debug.Stack()))
		}
	}()

	c.Send <- msg
}
