package ws

import (
	"PIM_Server/client"
	"PIM_Server/log"
	"PIM_Server/service"
	"PIM_Server/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	wsReadBufferSize  = 4096 // ws读最大buffer
	wsWriteBufferSize = 4096 // ws写最大buffer
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  wsReadBufferSize,
	WriteBufferSize: wsWriteBufferSize,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		log.Info("upgrade ua:", r.Header["User-Agent"], ", referer:", r.Header["Referer"])
		return true
	},
}

// onWebsocket ws链接事件
func onWebsocket(c *gin.Context) {
	log.Infof("new ws conn")
	conn, err := upgrade.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Errorf("ws upgrade err:", err)
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 用户连接事件
	// 首先根据token获取platform和uin
	cli := client.New(conn.RemoteAddr().String(), conn)
	cli.Uid = utils.GetUid(c)
	cli.PlatformId = client.PlatformWeb
	log.Infof("ws connected ok, uid:%d, platform:%d, conn:%s",
		cli.Uid, cli.PlatformId, conn.RemoteAddr().String())
	service.AddClient(cli)
}

// Routers websocket服务
func Routers(r *gin.Engine) {
	// websocket
	r.GET("/ws", onWebsocket)
}
