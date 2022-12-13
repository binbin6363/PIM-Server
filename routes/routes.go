package routers

import (
	"PIM_Server/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type Option func(*gin.Engine)

var options []Option

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.GetHeader("Origin") //请求头部
		if origin != "" {
			// 可将将* 替换为指定的域名
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 只有ws才校验token
		if c.Request.URL.Path != "/ws" {
			return
		}

		sc, err := service.DefaultService.ParseToken(c.Query("token"))
		if err != nil {
			log.Printf("parse token failed, disconnect ws, err:%+v", err)
			return
		}
		c.Set("uid", sc.Id)
		c.Set("username", sc.Audience)

		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

// Register 注册路由配置
func Register(opts ...Option) {
	options = append(options, opts...)
}

// Init 初始化
func Init() *gin.Engine {
	r := gin.Default()
	r.Use(Cors())
	r.Use(JWTAuthMiddleware())
	for _, opt := range options {
		opt(r)
	}
	return r
}
