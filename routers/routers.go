package routers

import (
	"PIM_Server/routers/auth"
	"PIM_Server/routers/chat"
	"PIM_Server/routers/ws"

	"github.com/gin-gonic/gin"
)

type Option func(*gin.Engine)

var options []Option

func init() {
	Register(auth.Routers)
	Register(chat.Routers)
	Register(ws.Routers)
}

// Register 注册路由配置
func Register(opts ...Option) {
	options = append(options, opts...)
}

func Routes() []Option {
	return options
}
