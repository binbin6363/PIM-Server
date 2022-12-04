package service

import (
	"PIM_Server/api"
	"PIM_Server/client"
	"encoding/json"
)

// OnTextMsg 文本消息通知本服务
func OnTextMsg(req *api.SendTextMsgReq) {
	login := &client.Login{
		PlatformID: 102,
		Uid:        20221113,
	}
	login.Client = clientManager.GetUserClient(login.PlatformID, login.Uid)

	byts, _ := json.Marshal(req)
	login.Client.SendMsg(byts)
}
