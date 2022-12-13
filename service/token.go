package service

import (
	"context"
	"log"
	"strconv"
	"time"

	"PIM_Server/config"
	"PIM_Server/model"
	"github.com/dgrijalva/jwt-go"
)

// CreateJwt 生成token信息
func (s *Service) CreateJwt(_ context.Context, userInfo *model.UserInfo) (error, string) {
	expiresTime := time.Now().Unix() + int64(config.AppConfig().ServerInfo.TokenExpire)
	claims := jwt.StandardClaims{
		Audience:  userInfo.UserName,               // 受众
		ExpiresAt: expiresTime,                     // 失效时间
		Id:        strconv.Itoa(int(userInfo.Uid)), // 编号
		IssuedAt:  time.Now().Unix(),               // 签发时间
		Issuer:    "pim",                           // 签发人
		NotBefore: time.Now().Unix(),               // 生效时间
		Subject:   "login",                         // 主题
	}
	var jwtSecret = []byte(config.AppConfig().ServerInfo.Secret)
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if token, err := tokenClaims.SignedString(jwtSecret); err == nil {
		log.Printf("gen token ok, %s(%d) token:[%s]", userInfo.UserName, userInfo.Uid, token)
		return nil, token
	} else {
		log.Printf("gen token failed, %s(%d) err:%v", userInfo.UserName, userInfo.Uid, err)
		return err, ""
	}
}

func (s *Service) ParseToken(token string) (*jwt.StandardClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (i interface{}, e error) {
		return []byte(config.AppConfig().ServerInfo.Secret), nil
	})
	if err == nil && jwtToken != nil {
		if claim, ok := jwtToken.Claims.(*jwt.StandardClaims); ok && jwtToken.Valid {
			return claim, nil
		}
	}

	log.Printf("token invalid, token:%s, jwtToken:%+v\n", token, jwtToken)
	return nil, err
}
