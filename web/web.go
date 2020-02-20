package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mywork/fundation"
	"mywork/sql"
	"net/http"
	"strconv"
)

func Registe(c *gin.Context)  {
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	telephone:=c.PostForm("telephone")
	if sql.UserSignup(username,password,telephone) {
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"注册成功"})
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"massage":"数据库Insert报错"})
	}
}
func Login(c *gin.Context)  {
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	telephone:=c.PostForm("telephone")
	fmt.Println(username,password,telephone)
	if username!=""{
		a,_:=sql.UserSignin(username,password)
		if a {
			c.SetCookie("username", username, 10, "/", "localhost", false, true)
			c.JSON(200,gin.H{"status":http.StatusOK,"message":"登录成功"})
		}else {
			c.JSON(403,gin.H{"status":http.StatusForbidden,"message":"登录失败，用户名或密码错误"})
		}
	}else {
		telephone="tele"+telephone
		a,username:=sql.UserSignin(telephone,password)
		if a {
			c.SetCookie("username", username, 10, "/", "localhost", false, true)
			c.JSON(200,gin.H{"status":http.StatusOK,"message":"登录成功"})
		}else {
			c.JSON(403,gin.H{"status":http.StatusForbidden,"message":"登录失败，用户名或密码错误"})
		}
	}

}
func GetComment(c *gin.Context,quid string)  {
	quid="comment_"+quid
	a:=sql.FindMessageByPid(0,quid)
	b:=fundation.JsonNested(a)
	c.JSON(200,b)
}
func Mapa(c *gin.Context)  {
	if cookie,err:=c.Request.Cookie("username");err==nil {
		name:=cookie.Value
		portrait:=sql.Perportrait(name)
		c.String(200,portrait)
		follower:=sql.Person(name,"sign","name")
		for _,n:=range follower {
			c.JSON(200,gin.H{"关注了":n})
		}
		fans:=sql.Person(name,"name","sign")
		for _,n:=range fans {
			c.JSON(200,gin.H{"粉丝":n})
		}
		questions:=sql.Ques(name)
		for _,n:=range questions  {
			c.JSON(200,gin.H{"问题":n})
		}
	}else {
		c.String(200,"你没有登录")
		fmt.Println(err)
	}
}
func Update(c *gin.Context) {
	var name string
	table:=c.Query("table")
	column:=c.Query("column")
	value:=c.PostForm("value")
	if cook,err:=c.Request.Cookie("username");err==nil{
		name=cook.Value
	}else {
		c.String(401,"你没有登录")
	}
	if a:=sql.Upinform(table,column,value,name);a{
		c.String(200,"修改成功")
	}
}
func Firpage(c *gin.Context)  {
	qunames,questions,pictures,keys,likenum:=sql.Firpagesort()
	c.String(200,"热榜")
	for m,_:=range qunames {
		c.JSON(200,gin.H{"问题名字":qunames[m],"问题描述":questions[m],"图片":pictures[m],"编号":keys[m],"点赞":likenum[m]})
	}
	qunames,questions,pictures,keys,likenum=sql.Firpageunsort()
	c.String(200,"推荐")
	for m,_:=range qunames {
		c.JSON(200,gin.H{"问题名字":qunames[m],"问题描述":questions[m],"图片":pictures[m],"编号":keys[m],"点赞":likenum[m]})
	}
}
func PutComment(c *gin.Context)  {
	if cookie,err:=c.Request.Cookie("username");err==nil{
		name:=cookie.Value
		comment:=c.PostForm("comment")
		id:=c.Query("quid")
		pid:=c.PostForm("id")
		fil,err:=c.FormFile("picture")
		if err!=nil {
			fmt.Println(err)
		}
		if fil!=nil {
			fundation.SaveUploadedFile(fil,"E:/GoProjects/src/mywork/templates/zhihu/comment/")
			sql.Upinform("comment_"+id,"picture","comment/"+fil.Filename,name)
		}
		if sql.PutMessage(id,name,comment,pid) {
			c.String(200,"发表成功")
		}else {
			c.String(500,"鬼知道怎么错了")
		}
	}else {
		c.String(401,"你没有登录")
	}
}
func Likequ(c *gin.Context)  {
	var name string
	if cookie,err:=c.Request.Cookie("username");err==nil {
		name=cookie.Value
	}else {
		c.String(401,"你没有登录")
	}
	if _,err:=c.Request.Cookie(name);err==nil {
		c.String(200,"不能重复点赞")
	}else {
		id:=c.Query("id")
		intid,err:=strconv.Atoi(id)
		if err!=nil {
			fmt.Println(err)
		}
		sql.Likequ(intid)
		c.SetCookie(name,"0",10,"/", "localhost", false, true)
	}
}
func Select(c *gin.Context)  {
	mess:=c.PostForm("mess")
	qunames,questions,pictures,keys,likenum:=sql.Select(mess)
	for m,_:=range qunames {
		c.JSON(200,gin.H{"问题名字":qunames[m],"问题描述":questions[m],"图片":pictures[m],"编号":keys[m],"点赞":likenum[m]})
	}
}
func Page(c *gin.Context)  {
	key:=c.Query("key")
	intkey,_:=strconv.Atoi(key)
	quname,question,picture,likenum:=sql.Prom(intkey)
	c.JSON(200,gin.H{"问题名字":quname,"问题描述":question,"图片":picture,"编号":key,"点赞":likenum})
	GetComment(c,key)

}
func PutQuestion(c *gin.Context)  {
	if cookie,err:=c.Request.Cookie("username");err==nil {
		name:=cookie.Value
		question:=c.PostForm("question")
		quname:=c.PostForm("quname")
		fil,err:=c.FormFile("picture")
		if err!=nil {
			fmt.Println(err)
		}
		if fil!=nil {
			fundation.SaveUploadedFile(fil,"E:/GoProjects/src/mywork/templates/zhihu/comment/")
			sql.Upinform("question","picture","question/"+fil.Filename,name)
		}
		if sql.PutQuestion(quname,name,question) {
			 c.String(200,"发表成功")
		}else {
		 	c.String(500,"鬼知道怎么错了")
		}
	}else {
		c.String(401,"你没有登录")
	}
}