package api

type SendTextMsgNotice struct {
	Event   string             `json:"event"`
	Content SendTextMsgContent `json:"content"`
}

type SendTextMsgEvtRsp struct {
	Event   string             `json:"event"`
	Content SendTextMsgContent `json:"content"`
}

type SendTextMsgContent struct {
	Data       SendTextMsgData `json:"data"`
	ReceiverId int64           `json:"receiver_id"`
	SenderId   int64           `json:"sender_id"`
	TalkType   int             `json:"talk_type"`
}

type SendTextMsgData struct {
	Id         int64  `json:"id"`
	Sequence   int64  `json:"sequence"`
	TalkType   int    `json:"talk_type"`
	MsgType    int    `json:"msg_type"`
	UserId     int64  `json:"user_id"`
	ReceiverId int64  `json:"receiver_id"`
	Nickname   string `json:"nickname"`
	Avatar     string `json:"avatar"`
	IsRevoke   int    `json:"is_revoke"`
	IsMark     int    `json:"is_mark"`
	IsRead     int    `json:"is_read"`
	Content    string `json:"content"`
	CreatedAt  string `json:"created_at"`
}

// CommRsp 针对请求的ack
type CommRsp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ClientEventReq 客户端上行的消息
type ClientEventReq struct {
	Event   string      `json:"event"`
	Content string      `json:"content"`
	Data    interface{} `json:"data"`
}

// ClientEventRsp 客户端下行的消息
type ClientEventRsp struct {
	Event   string      `json:"event"`
	Content string      `json:"content"`
	Data    interface{} `json:"data"`
}
