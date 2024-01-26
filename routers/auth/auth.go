package auth

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func loginHandler(c *gin.Context) {
	//service.OnLogin()
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "success",
		"data":    nil,
	})
}

// Routers .
func Routers(r *gin.Engine) {
	login := r.Group("/notice/auth")
	{
		login.POST("/login", loginHandler)
	}
}
