package chat

import (
	"PIM_Server/api"
	"PIM_Server/log"
	"PIM_Server/service"
	"PIM_Server/utils"
	"github.com/gin-gonic/gin"
)

func onTextMsgHandler(c *gin.Context) {

	notice := &api.SendTextMsgNotice{}
	if err := c.ShouldBind(notice); err != nil {
		log.Infof("bind text msg notice failed, err:%+v", err)
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    4000,
			Message: err.Error(),
		})
		return
	}

	err := service.DefaultService.OnTextMsg(c, notice)
	if err != nil {
		utils.SendJsonRsp(c, &api.CommRsp{
			Code:    1000,
			Message: err.Error(),
		})
	} else {
		utils.SendJsonRsp(c, nil)
	}
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/notice/message")
	{
		login.POST("/text", onTextMsgHandler)
	}
}
