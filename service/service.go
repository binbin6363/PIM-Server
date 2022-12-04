package service

import (
	"PIM_Server/client"
)

const (
	defaultAppId = 101 // 默认平台Id
)

var (
	clientManager = client.NewManager()                   // 管理者
	platformIds   = []uint32{defaultAppId, 102, 103, 104} // 全部的平台

	serverIp   string
	serverPort string
)

// StartClientMgr 启动client管理服务
func StartClientMgr() {
	// 启动websocket客户端管理器
	go clientManager.Start()

}
