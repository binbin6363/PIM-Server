package service

import (
	"PIM_Server/client"
	"PIM_Server/config"
	"PIM_Server/dao"
)

const (
	defaultAppId = 101 // 默认平台Id
)

var (
	clientManager = client.NewManager()                   // 管理者
	platformIds   = []uint32{defaultAppId, 102, 103, 104} // 全部的平台

	serverIp       string
	serverPort     string
	DefaultService *Service
)

// StartClientMgr 启动client管理服务
func StartClientMgr() {
	// 启动websocket客户端管理器
	go clientManager.Start()
}

func AddClient(cli *client.Client) {
	clientManager.Register <- cli
}

// Service is service logic object
type Service struct {
	dao *dao.Dao
}

// New creates service instance
func New() *Service {
	srv := Service{
		dao: dao.New(config.AppConfig().DBInfo.Dsn,
			config.AppConfig().ServerInfo.DataCenterId,
			config.AppConfig().ServerInfo.WorkerId),
	}
	return &srv
}

// Init .
func Init() {
	DefaultService = New()
}
