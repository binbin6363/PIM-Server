package client

import (
	"PIM_Server/api"
	"PIM_Server/log"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	"io"
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
		log.Errorf("conn close failed, addr:%s, platform:%d,uid:%d,conn time:%d,heartbeat time:%d,login time:%d",
			c.Addr, c.PlatformId, c.Uid, c.ConnTime, c.HeartbeatTime, c.LoginTime)
	} else {
		log.Errorf("conn close success, addr:%s, platform:%d,uid:%d,conn time:%d,heartbeat time:%d,login time:%d",
			c.Addr, c.PlatformId, c.Uid, c.ConnTime, c.HeartbeatTime, c.LoginTime)
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

func (c *Client) ReceiveMessage() {
	for {
		t, d, e := c.Socket.ReadMessage()
		if e == io.EOF || e == io.ErrUnexpectedEOF {
			// close
			log.Errorf("conn closed by peer, uid:%d", c.Uid)
			c.Close()
			return
		} else if e != nil {
			log.Errorf("read conn err: %v, uid:%d", e, c.Uid)
			return
		}

		// handle message
		c.HandleMessage(t, d)
	}
}

func (c *Client) HandleMessage(messageType int, data []byte) {
	switch messageType {
	case websocket.TextMessage:
		c.HandleText(data)
	case websocket.BinaryMessage:
		c.HandleBinary(data)
	case websocket.CloseMessage:
		c.Close()
	case websocket.PingMessage:
		c.HandlePing(data)
	case websocket.PongMessage:
		// nothing
	}
}

func (c *Client) HandleText(data []byte) {
	log.Infof("receive text from uid:%d", c.Uid)
	msg := &api.ClientEventReq{}
	if err := json.Unmarshal(data, msg); err != nil {
		log.Errorf("parse client msg err, err:%v", err)
		return
	}

	rsp := &api.ClientEventRsp{
		Event: msg.Event,
	}
	c.HeartbeatTime = time.Now().Unix()
	c.SendJsonMsg(rsp)
}

func (c *Client) HandleBinary(data []byte) {
	log.Infof("receive binary from uid:%d", c.Uid)
}

func (c *Client) HandlePing(data []byte) {
	log.Infof("receive ping from uid:%d", c.Uid)
	c.HeartbeatTime = time.Now().Unix()
	c.SendJsonMsg(map[string]string{
		"event": "heartbeat",
		"time":  cast.ToString(c.HeartbeatTime),
	})
}

func (c *Client) Run() {
	go c.ReceiveMessage()
	for {
		select {
		case data, ok := <-c.Send:
			if !ok {
				log.Errorf("conn closed, exit. uid::%d", c.Uid)
				c.Send = nil
				return
			}
			// 发送数据
			log.Infof("real send data to socket")
			if err := c.Socket.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Errorf("write message failed, uid:%d", c.Uid)
			}

		}
	}
}
