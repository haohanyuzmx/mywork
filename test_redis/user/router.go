package user

import (
	"github.com/gin-gonic/gin"
	"test_redis/middleware"
)

func User(u *gin.Engine)  {
	u.POST("register",Registe)
	u.GET("charts",Charts)
	u.POST("login",Login)
	u.POST("vote",middleware.Inspect,Vote)
	u.GET("matchon",middleware.Inspect,Matchon)
	u.GET("metchoff",middleware.Inspect,Matchoff)
}