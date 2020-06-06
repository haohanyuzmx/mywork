package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	fundation "gobang"
	"net/http"
)

type Room struct {
	Roomid  string
	Member1 string
	Member2 string
}

func AddRoom(ctx *gin.Context) {
	var goroom Room
	goroom.Roomid = ctx.Query("roomid")
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
	//var u fundation.User
	conn.Do("sadd", goroom.Roomid, member)
	ctx.SetCookie("roomid", goroom.Roomid, 10, "/", "localhost", false, true)
	ctx.JSON(200,gin.H{
		"code":http.StatusMovedPermanently,
		"mess":"请用ws连接/chess",
	})
}
