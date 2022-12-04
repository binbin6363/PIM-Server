package chat

import (
	"PIM_Server/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func textMsgHandler(c *gin.Context) {

	rsp := api.SendTextMsgRsp{
		Id:         10001,
		TalkType:   1,
		ReceiverId: 20221114,
		Name:       "20221114",
		RemarkName: "mark",
		Avatar:     "",
		IsDisturb:  0,
		IsTop:      0,
		IsOnline:   0,
		IsRobot:    0,
		UnreadNum:  1,
		Content:    "chat content",
		DraftText:  "",
		MsgText:    "chat MsgText",
		IndexName:  "",
		CreatedAt:  "20221119",
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    rsp,
	})

	req := &api.SendTextMsgReq{}
	c.BindJSON(req)
	//service.OnTextMsg(req)
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/notice/message")
	{
		login.POST("/text", textMsgHandler)
	}
}
