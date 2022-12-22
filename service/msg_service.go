package service

import (
	"PIM_Server/api"
	"PIM_Server/client"
	"PIM_Server/log"
	"context"
)

// OnTextMsg 文本消息通知本服务
func (s *Service) OnTextMsg(ctx context.Context, req *api.SendTextMsgNotice) (err error) {
	log.Infof("enter OnTextMsg %+v", req)

	if req.Content.TalkType == 1 {
		receiver := clientManager.GetUserClient(client.PlatformWeb, req.Content.ReceiverId)

		if receiver != nil {
			req.Content.Data.UserId = req.Content.SenderId
			err = receiver.SendJsonMsg(req)
		} else {
			log.Infof("user receiver %d not connected", req.Content.ReceiverId)
		}
	} else if req.Content.TalkType == 2 {
		// 获取群成员列表
		groupId := req.Content.ReceiverId
		_, groupMembers := s.dao.GetGroupMemberList(ctx, groupId)
		for _, groupMember := range groupMembers {
			if groupMember.Uid == req.Content.SenderId {
				log.Debugf("skip notice self:%d", req.Content.SenderId)
				continue
			}
			// 获取链接，推送
			receiver := clientManager.GetUserClient(client.PlatformWeb, groupMember.Uid)

			if receiver != nil {
				req.Content.Data.UserId = req.Content.SenderId
				err = receiver.SendJsonMsg(req)
				log.Infof("notice user uid:%d", groupMember.Uid)
			} else {
				log.Infof("user receiver %d not connected", groupMember.Uid)
			}
		}
	}

	return err
}
