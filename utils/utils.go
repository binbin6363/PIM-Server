package utils

import (
	"PIM_Server/api"
	"PIM_Server/log"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// EncryptPassword 对密码加密
func EncryptPassword(password string) (error, string) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return err, string(bytes)
}

// CheckPasswordHash 密码校验
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUid 从上下文获取uid
func GetUid(c *gin.Context) (uid int64) {
	if aUid, ok := c.Get("uid"); ok {
		uid = cast.ToInt64(aUid)
		return uid
	}
	return 0
}

// GetUserName 从上下文获取username
func GetUserName(c *gin.Context) (err error, username string) {
	if aUserName, ok := c.Get("username"); ok {
		username = cast.ToString(aUserName)
		return nil, username
	}
	return errors.New("gin get username failed"), ""
}

// GetInt64Val 从上下文获取val
func GetInt64Val(c *gin.Context, key string) (val int64) {
	if v, ok := c.Get(key); ok {
		val = cast.ToInt64(v)
		return val
	}
	return 0
}

// GetString 从上下文获取val
func GetString(c *gin.Context, key string) (val string) {
	if v, ok := c.Get(key); ok {
		val = cast.ToString(v)
		return val
	}
	return ""
}

// QueryInt64Val 从上下文获取val
func QueryInt64Val(c *gin.Context, key string) (val int64) {
	val = cast.ToInt64(c.Query(key))
	return val
}

// QueryString 从上下文获取val
func QueryString(c *gin.Context, key string) (val string) {
	val = c.Query(key)
	return val
}

// SendJsonRsp 回复json消息
func SendJsonRsp(c *gin.Context, rsp *api.CommRsp) {
	if rsp == nil {
		log.InfoContextf(c, "handle ok, send empty rsp")
		c.JSON(http.StatusOK, gin.H{
			"code":    200,
			"message": "ok",
			"data":    nil,
		})
	} else if rsp.Code == 0 || rsp.Code == 200 {
		rsp.Code = 200
		log.InfoContextf(c, "handle ok, send rsp:%+v", rsp)
		c.JSON(http.StatusOK, rsp)
	} else {
		log.InfoContextf(c, "handle error, code:%d, message:%s", rsp.Code, rsp.Message)
		c.JSON(http.StatusOK, rsp)
	}
}

// FormatTimeStr 格式化时间戳到字符串
func FormatTimeStr(nowTime int64) string {
	return time.Unix(nowTime, 0).Format("2006-01-02 15:04:05")
}

// SendJsonRsp 回复json消息
//func SendJsonRsp(c *gin.Context, rsp *api.CommRsp) {
//	if rsp == nil || rsp.Code == 0 {
//		log.Infof("handle ok, send rsp")
//		c.JSON(http.StatusOK, gin.H{
//			"code":    0,
//			"message": "ok",
//			"data":    nil,
//		})
//	} else {
//		log.Infof("handle error, code:%d, message:%s", rsp.Code, rsp.Message)
//		c.JSON(http.StatusOK, rsp)
//	}
//}
