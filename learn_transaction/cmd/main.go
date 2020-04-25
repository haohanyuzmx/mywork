package main

import (
	"github.com/gin-gonic/gin"
	"learn_transaction/user"
)

func main()  {
	r:=gin.Default()
	r.POST("chose",user.Chose)
	r.Run()
}
