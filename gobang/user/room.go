package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	fundation "gobang"
)

type Room struct {
	Roomid  string
	Member1 string
	Member2 string
}


//AddRoom从body中获取roomid,在redis中建立房间
func AddRoom(ctx *gin.Context) {
	var goroom Room
	goroom.Roomid = ctx.PostForm("roomid")
	//fmt.Println(goroom.Roomid)
	conn := fundation.Pool.Get()
	mem, err := conn.Do("scard", goroom.Roomid)
	if wrongSend(ctx, err, "redis错误") {
		return
	}
	intmem := mem.(int64)
	if intmem >= 2 {
		err := errors.New("成员大于2")
		wrongSend(ctx, err, "房间成员大于2")
		return
	}
	member, err := ctx.Cookie("username")
	if wrongSend(ctx,err,"没有登录") {
		return
	}
	//var u fundation.User
	conn.Do("sadd", goroom.Roomid, member)
	ctx.SetCookie("roomid", goroom.Roomid, 300, "/", "localhost", false, true)
	ctx.JSON(200,gin.H{
		"code":003,
		"mess":"请用ws连接/chess",
	})
}
