package model

// UserInfo 用户信息表
type UserInfo struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid        int64  `gorm:"column:uid"`
	UserName   string `gorm:"column:username"`
	PassHash   string `gorm:"column:passhash"`
	NickName   string `gorm:"column:nickname"`
	Phone      string `gorm:"column:phone"`
	Email      string `gorm:"column:email"`
	Avatar     string `gorm:"column:avatar"`
	Gender     int    `gorm:"column:gender"`
	UserTag    string `gorm:"column:user_tag"`
	Motto      string `gorm:"column:motto"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (UserInfo) TableName() string {
	return "user_infos"
}

// Contacts 联系人好友
type Contacts struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid        int64  `gorm:"column:uid"`
	ContactID  int64  `gorm:"column:contact_id"`
	RemarkName string `gorm:"column:remark_name"`
	Status     int    `gorm:"column:status"`
	Sequence   int64  `gorm:"column:sequence"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Contacts) TableName() string {
	return "contacts"
}

// UserContact 好友关系，联表查询结果
type UserContact struct {
	Uid          int64  `gorm:"column:uid"`
	Nickname     string `gorm:"column:nickname"`
	Gender       int    `gorm:"column:gender"`
	Motto        string `gorm:"column:motto"`
	Avatar       string `gorm:"column:avatar"`
	FriendRemark string `gorm:"column:remark_name"`
}

// Conversations 会话列表
type Conversations struct {
	ID                 int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid                int64  `gorm:"column:uid"`
	ContactID          int64  `gorm:"column:contact_id"`
	ConversationType   int    `gorm:"column:conversation_type"`
	ConversationName   string `gorm:"column:conversation_name"`
	ConversationStatus int    `gorm:"column:conversation_status"`
	Unread             int    `gorm:"column:unread"`
	MsgDigest          string `gorm:"column:msg_digest"`
	MsgID              int64  `gorm:"column:msgid"`
	CreateTime         int64  `gorm:"column:create_time"`
	UpdateTime         int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (Conversations) TableName() string {
	return "conversations"
}

// SingleMessages 单人消息
type SingleMessages struct {
	ID          int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	Uid         int64  `gorm:"column:uid"`
	MsgID       int64  `gorm:"column:msgid"`
	ClientMsgID int64  `gorm:"column:client_msgid"`
	SenderID    int64  `gorm:"column:senderid"`
	ReceiverID  int64  `gorm:"column:receiverid"`
	MsgType     int    `gorm:"column:msg_type"`
	Content     string `gorm:"column:content"`
	MsgStatus   int    `gorm:"column:msg_status"`
	CreateTime  int64  `gorm:"column:create_time"`
	UpdateTime  int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (SingleMessages) TableName() string {
	return "single_messages"
}

// GroupMessages 群消息
type GroupMessages struct {
	ID         int64  `gorm:"column:id;primarykey;AUTO_INCREMENT"`
	GroupID    int64  `gorm:"column:groupid"`
	MsgID      int64  `gorm:"column:msgid"`
	SenderID   int64  `gorm:"column:senderid"`
	MsgType    int    `gorm:"column:msg_type"`
	Content    string `gorm:"column:content"`
	MsgStatus  int    `gorm:"column:msg_status"`
	CreateTime int64  `gorm:"column:create_time"`
	UpdateTime int64  `gorm:"column:update_time"`
}

// TableName 默认是通过结构体的蛇形复数来指定表名，这里通过TableName显示定义出来，便于问题排查
func (GroupMessages) TableName() string {
	return "group_messages"
}
