package service

import (
	"PIM_Server/api"
	"PIM_Server/client"
	"log"
)

// OnTextMsg 文本消息通知本服务
func (s *Service) OnTextMsg(req *api.SendTextMsgNotice) (err error) {
	log.Printf("[DEBUG] enter OnTextMsg %+v\n", req)

	//sender := clientManager.GetUserClient(client.PlatformWeb, req.Content.SenderId)
	receiver := clientManager.GetUserClient(client.PlatformWeb, req.Content.ReceiverId)
	//if sender != nil {
	//	req.Content.Data.UserId = req.Content.SenderId
	//	err = sender.SendJsonMsg(req)
	//} else {
	//	log.Printf("[ERROR] user sender %d not connected\n", req.Content.SenderId)
	//}
	if receiver != nil {
		req.Content.Data.UserId = req.Content.SenderId
		err = receiver.SendJsonMsg(req)
	} else {
		log.Printf("[NOTICE] user receiver %d not connected\n", req.Content.ReceiverId)
	}

	return err
}
