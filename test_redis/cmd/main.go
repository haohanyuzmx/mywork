package main

import (
	"github.com/gin-gonic/gin"
	"test_redis/user"
)

func main()  {
	r:=gin.Default()
	user.User(r)
	r.Run()
}