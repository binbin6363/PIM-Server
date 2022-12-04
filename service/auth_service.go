package service

import "PIM_Server/client"

// OnLogin 登录通知本服务
func OnLogin() {
	login := &client.Login{
		PlatformID: 102,
		Uid:        20221113,
	}
	login.Client = clientManager.GetUserClient(login.PlatformID, login.Uid)

	clientManager.Login <- login
}

// OnLogout 登出通知本服务
func OnLogout() {
	//clientManager.Login <- nil
}
