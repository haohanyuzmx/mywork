package main

import (
	"github.com/gin-gonic/gin"
	"gobang/user"
)

func main()  {

	r:=gin.Default()
	user.Chess(r)
	r.Run()
}
