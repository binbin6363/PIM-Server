package utils

import (
	"PIM_Server/api"
	"PIM_Server/log"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

// EncryptPassword 对密码加密
func EncryptPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckPasswordHash 密码校验
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetUid 从上下文获取uid
func GetUid(c *gin.Context) (err error, uid int64) {
	if aUid, ok := c.Get("uid"); ok {
		uid = cast.ToInt64(aUid)
		return nil, uid
	}
	return errors.New("gin get uid failed"), 0
}

// GetUserName 从上下文获取username
func GetUserName(c *gin.Context) (err error, username string) {
	if aUserName, ok := c.Get("username"); ok {
		username = cast.ToString(aUserName)
		return nil, username
	}
	return errors.New("gin get username failed"), ""
}

// SendJsonRsp 回复json消息
func SendJsonRsp(c *gin.Context, rsp *api.CommRsp) {
	if rsp == nil || rsp.Code == 0 {
		log.Infof("handle ok, send rsp")
		c.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "ok",
			"data":    nil,
		})
	} else {
		log.Infof("handle error, code:%d, message:%s", rsp.Code, rsp.Message)
		c.JSON(http.StatusOK, rsp)
	}
}
