package service

import (
	"PIM_Server/client"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		log.Println("upgrade ua:", r.Header["User-Agent"], ", referer:", r.Header["Referer"])
		return true
	},
}

// Websocket 服务
func Websocket(r *gin.Engine) {
	// websocket
	r.GET("/ws", func(context *gin.Context) {
		c, err := upgrade.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			log.Print("upgrade:", err)
			http.NotFound(context.Writer, context.Request)
			return
		}

		log.Println("webSocket conn:", c.RemoteAddr().String())

		// 用户连接事件
		// 首先根据token获取platform和uin
		cli := client.New(c.RemoteAddr().String(), c)
		cli.Uid = 20221113
		cli.PlatformId = 102
		clientManager.Register <- cli
	})
}
