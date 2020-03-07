package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	fundation "test_1"
	"test_1/jwt"
)

func Inspect(c *gin.Context)  {
	auth:=c.GetHeader("Authorization")
	if len(auth)<7 {
		c.JSON(http.StatusOK,gin.H{
			"code":300,
			"mess":"登录失败",
		})
		c.Abort()
		return
	}
	token:=auth[7:]
	username,err:=jwt.Checktoken(token)
	if !fundation.Wrong(err){
		c.Abort()
		return
	}
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"mess":"认证成功",
	})
	c.Set("username",username)
	c.Next()
}
