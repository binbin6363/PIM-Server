package plugins

import (
	"math/rand"
	"net/http"
	"time"

	"PIM_Server/config"
	"PIM_Server/log"
	"PIM_Server/routers"
	"PIM_Server/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	gindump "github.com/tpkeeper/gin-dump"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.uber.org/zap"
)

// Init gin初始化
func Init(serviceName string) *gin.Engine {

	r := gin.Default()
	r.Use(Cors())
	r.Use(JWTAuthMiddleware())
	r.Use(otelgin.Middleware(serviceName))
	r.Use(ZapTraceLogger())
	if config.AppConfig().ServerInfo.DebugReqRsp {
		log.Infof("open req/rsp debug log")
		r.Use(gindump.DumpWithOptions(true, true, true, false, false, func(dumpStr string) {
			log.Infof("dump: [%s]", dumpStr)
		}))
	}
	log.Infof("show routers: %v", routers.Routes())
	for _, opt := range routers.Routes() {
		opt(r)
	}

	rand.Seed(time.Now().UnixNano())
	return r
}

// Cors 跨域处理中间件
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
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 只有ws才校验token
		if c.Request.URL.Path != "/ws" {
			return
		}

		sc, err := service.DefaultService.ParseToken(c.Query("token"))
		if err != nil {
			log.Infof("parse token failed, disconnect ws, err:%+v", err)
			return
		}
		c.Set("uid", sc.Id)
		c.Set("username", sc.Audience)

		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}

// ZapTraceLogger 创建一个ZapTraceLogger中间件
func ZapTraceLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Gin Context 中获取 Trace ID（假设 Trace ID 存储在 Header 中）
		traceID := c.Request.Header.Get(log.LoggerTraceID)
		if len(traceID) == 0 {
			traceID = cast.ToString(rand.Uint64())
		}

		// 将 Trace ID 添加到 Zap Logger 的上下文字段中
		loggerWithTraceID := log.GetLogger().With(zap.String(log.LoggerTraceID, traceID))

		// 将 Zap Logger 添加到 Gin Context 中，以便在请求处理程序中使用
		c.Set(log.LoggerTag, loggerWithTraceID)

		// 继续处理请求
		t := time.Now()
		log.InfoContextf(c, "recv msg")
		c.Next()
		log.InfoContextf(c, "done, cost: %v", time.Since(t))
	}
}
