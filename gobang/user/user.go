package user

import (
	"errors"
	"github.com/gin-gonic/gin"
	fundation "gobang"
	"log"
)

type Userinfor struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
//Login通过json读取信息并查询数据库，若成功设置cookie保存登录状态
func Login(ctx *gin.Context) {
	var uinf Userinfor
	err := ctx.ShouldBindJSON(&uinf)
	if wrongSend(ctx, err, "数据不全") {
		return
	}
	var u fundation.User
	err = fundation.DB.Where(&fundation.User{Username: uinf.Username, Password: uinf.Password}).First(&u).Error
	if u.ID == 0 {
		wrongSend(ctx,err,"账号密码出错")
		return
	}
	ctx.SetCookie("username", u.Username, 600, "/", "localhost", false, true)
	ctx.JSON(200, gin.H{
		"code": 200,
		"mess": "成功",
	})
}
//Register通过json读取信息并查询数据库，不冲突则写入数据库
func Register(ctx *gin.Context) {
	var uinf Userinfor
	err := ctx.ShouldBindJSON(&uinf)
	if wrongSend(ctx, err, "数据不全") {
		return
	}
	var u fundation.User
	if !fundation.DB.HasTable(&u){
		fundation.DB.CreateTable(&u)
	}
	u.Username = uinf.Username
	fundation.DB.Table("users").Where("username=?", uinf.Username).First(&u)
	if u.ID != 0 {
		err=errors.New("账号重复")
		wrongSend(ctx,err,"账号重复")
		return
	}
	u.Password = uinf.Password
	fundation.DB.Create(&u)
	ctx.JSON(200, gin.H{
		"code": 200,
		"mess": "成功",
	})
}

func wrongSend(ctx *gin.Context, err error, mess string) bool {
	if err != nil {
		ctx.JSON(200, gin.H{
			"Code": 001,
			"mess": mess,
		})
		log.Println(err)
		return true
	}
	return false
}
