package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	jwt2 "test_redis/jwt"
)

func Inspect(c *gin.Context) {
	auth:=c.GetHeader("Authorization")
	if len(auth)<7 {
		c.JSON(200,gin.H{
			"code":500,
			"mess":"请求头错误",
		})
		c.Abort()
		return
	}
	token:=auth[7:]
	var jwt jwt2.JWT
	if jwt.Check(token)!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"mess":"token伪造",
		})
		c.Abort()
		return
	}
	c.Set("username",jwt.Payload.Username)
	fmt.Println(jwt.Payload)
	c.Next()
}
