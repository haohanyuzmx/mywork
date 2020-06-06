package user

import "github.com/gin-gonic/gin"

func User(u *gin.Engine)  {
	userfunc:=u.Group("user")
	userfunc.POST("register",Register)
	userfunc.POST("login",Login)
}
func Chess(u *gin.Engine)  {
	go Manager.Start()
	chessnow:=u.Group("chess")
	chessnow.POST("addroom",AddRoom)
	chessnow.GET("room",TestHandler)
}
