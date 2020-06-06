package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	fundation "gobang"
)

type Userinfor struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(ctx *gin.Context) {
	var uinf Userinfor
	err := ctx.ShouldBindJSON(&uinf)
	if wrongSend(ctx, err, "数据不全") {
		return
	}
	var u fundation.User
	err = fundation.DB.Where(&fundation.User{Username: uinf.Username, Password: uinf.Password}).First(&u).Error
	if u.ID == 0 {
		ctx.JSON(200, gin.H{
			"code": 500,
			"mess": "错误",
		})
		fmt.Println(err)
		return
	}
	ctx.SetCookie("username", u.Username, 10, "/", "localhost", false, true)
}
func Register(ctx *gin.Context) {
	var uinf Userinfor
	err := ctx.ShouldBindJSON(&uinf)
	if wrongSend(ctx, err, "数据不全") {
		return
	}
	var u fundation.User
	u.Username = uinf.Username
	fundation.DB.Where("username", uinf.Username).First(&u)
	if u.ID != 0 {
		ctx.JSON(200, gin.H{
			"Code": 500,
			"mess": "账号重合",
		})
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
			"Code": 500,
			"mess": mess,
		})
		fmt.Println(err)
		return true
	}
	return false
}
