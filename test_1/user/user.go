package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	fundation "test_1"
	"test_1/jwt"
)

type Userinfor struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
func Register(c *gin.Context)  {
	var i Userinfor
	var u fundation.User
	if err:=c.ShouldBindJSON(&i);err!=nil{
	c.JSON(http.StatusOK,gin.H{
		"code":400,
		"mess":"数据不够",
	})
		return
	}
	u.Username=i.Username
	u.Password=i.Password
	fundation.DB.Create(&u)
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"mess":"注册成功",
	})
}
func Login(c *gin.Context)  {
	var i Userinfor
	var u fundation.User
	c.ShouldBindJSON(&i)
	fmt.Println(i)
	fundation.DB.Where("username=? and password=?",i.Username,i.Password).First(&u)
	if u.ID==0 {
		c.JSON(http.StatusOK,gin.H{
			"code":300,
			"mess":"密码或者账号错误",
		})
		return
	}
	token:=jwt.Creat(u.Username)
	c.JSON(http.StatusOK,gin.H{
		"code":200,
		"token":token,
	})
}
func All(c *gin.Context)  {
	var us []fundation.User
	var json gin.H
	var jsons []gin.H
	fundation.DB.Select("*").Find(&us)
	for _,n:=range us {
		json=gin.H{
			"username":n.Username,
		}
		jsons=append(jsons,json)
	}
	c.JSON(http.StatusOK,jsons)
}
func Update(c *gin.Context)  {
	username:=c.Keys["username"]
	column:=c.PostForm("column")
	value:=c.PostForm("value")
	fmt.Println(username)
	err:=fundation.DB.Table("users").Where("username=?",username).Update(column,value).Error
	if !fundation.Wrong(err){
	 	c.JSON(http.StatusOK,gin.H{
	 		"code":300,
	 		"mess":"错误",
		})
		return
	}
	 c.JSON(http.StatusOK,gin.H{
	 	"code":200,
	 	"mess":"成功",
	 })
}