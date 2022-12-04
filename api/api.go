package api

type SendTextMsgReq struct {
	ReceiverId int    `json:"receiver_id"`
	TalkType   int    `json:"talk_type"`
	Text       string `json:"text"`
}

type SendTextMsgRsp struct {
	Id         int    `json:"id"`
	TalkType   int    `json:"talk_type"`
	ReceiverId int    `json:"receiver_id"`
	Name       string `json:"name"`
	RemarkName string `json:"remark_name"`
	Avatar     string `json:"avatar"`
	IsDisturb  int    `json:"is_disturb"`
	IsTop      int    `json:"is_top"`
	IsOnline   int    `json:"is_online"`
	IsRobot    int    `json:"is_robot"`
	UnreadNum  int    `json:"unread_num"`
	Content    string `json:"content"`
	DraftText  string `json:"draft_text"`
	MsgText    string `json:"msg_text"`
	IndexName  string `json:"index_name"`
	CreatedAt  string `json:"created_at"`
}
