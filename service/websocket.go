package service

import (
	"PIM_Server/client"
	"PIM_Server/log"
	"PIM_Server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		log.Info("upgrade ua:", r.Header["User-Agent"], ", referer:", r.Header["Referer"])
		return true
	},
}

// Websocket 服务
func Websocket(r *gin.Engine) {
	// websocket
	r.GET("/ws", func(context *gin.Context) {
		log.Infof("new conn")
		c, err := upgrade.Upgrade(context.Writer, context.Request, nil)
		if err != nil {
			log.Errorf("upgrade:", err)
			http.NotFound(context.Writer, context.Request)
			return
		}

		// 用户连接事件
		// 首先根据token获取platform和uin
		cli := client.New(c.RemoteAddr().String(), c)

		if err, uid := utils.GetUid(context); err == nil {
			cli.Uid = uid
		} else {
			log.Infof("get uid from token failed, err:%+v", err)
			return
		}
		cli.PlatformId = client.PlatformWeb
		log.Infof("ws connected ok, uid:%d, platform:%d, conn:%s",
			cli.Uid, cli.PlatformId, c.RemoteAddr().String())
		clientManager.Register <- cli
	})
}
