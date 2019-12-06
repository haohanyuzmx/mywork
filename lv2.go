package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
)
var db *sql.DB

func init()  {
	db,_=sql.Open("mysql","root:@tcp(localhost:3306)/chat?charset=utf8")
	db.SetMaxOpenConns(1000)
	err:=db.Ping()
	if err!=nil {
		fmt.Println("fail to connnect to db")
		fmt.Println(err.Error())
	}
}
func DBconn() *sql.DB {
	return db
}
func UserSignup(username string,passward string) bool {
	stmt,err:=DBconn().Prepare("insert into user(username,password) value (?,?)")
	if err!=nil{
		fmt.Println("fail to insert1")
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_,err=stmt.Exec(username,passward)
	if err!=nil{
		fmt.Println("fail to insert2")
		fmt.Println(err.Error())
		return false
	}
	return true
}
func UserSignin(username string,passward string) bool {
	var id int
	stmt,err:=DBconn().Prepare("select id from user where username=? and password=?")
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	err=stmt.QueryRow(username,passward).Scan(&id)
	if err!=nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}
func SendMesage(username,message string,uid int) bool {
	stmt,err:=DBconn().Prepare("insert into comment(name,message,uid) value (?,?,?)")
	if err!=nil{
		fmt.Println("fail to insert1")
		fmt.Println(err.Error())
		return false
	}
	_,err=stmt.Exec(username,message,uid)
	if err!=nil{
		fmt.Println("fail to insert2")
		fmt.Println(err.Error())
		return false
	}
	return true
}
func GetMesage(i int) ( map[int]map[string]string) {
	comment:=make(map[string]string)
	var id,uid,test,likenum int=0,0,0,0
	var username,message string
	floor:=make(map[int]map[string]string)
	for a:=(i-1)*10;a<=i*10 ;a++  {
		stmt,err:=DBconn().Prepare("SELECT * FROM comment limit ?,?")
		if err!=nil {
			fmt.Println(err)
		}
		err=stmt.QueryRow(a,a+1).Scan(&id,&username,&message,&test,&likenum)
		username="id为:"+strconv.Itoa(id)+" 用户名:"+username
		message=message+" 点赞数为:"+strconv.Itoa(likenum)
		if err!=nil {
			fmt.Println(err)
		}else {
			if test!=uid {
				uid=test
				comment:=make(map[string]string)
				comment[username]=message
				floor[uid]=comment
			}else {
				comment[username]=message
				floor[uid]=comment
			}
		}
	}
	return floor
}
func query(name string) bool {
	var id int
	stmt,err:=DBconn().Prepare("select id from user where username=?")
	if err!=nil {
		fmt.Println(err)
		return false
	}
	err=stmt.QueryRow(name).Scan(&id)
	if err!=nil {
		fmt.Println(err)
		return false
	}
	return true
}
func like(a string)  {
	var num int
	stmt,err:=DBconn().Prepare("SELECT likenum FROM comment where id=?")
	if err!=nil {
		fmt.Println(err)
	}
	i,err:=strconv.Atoi(a)
	err=stmt.QueryRow(i).Scan(&num)
	num++
	if err!=nil {
		fmt.Println(err)
	}
	stmt,err=DBconn().Prepare("update comment set likenum=? where id=?")
	if err!=nil {
		fmt.Println(err)
	}
	_,err=stmt.Exec(num,i)
	if err!=nil {
		fmt.Println(err)
	}
}
func affirm(upid int,name string) bool {
	var id int
	stmt,err:=DBconn().Prepare("select id from comment where name=?")
	if err!=nil {
		fmt.Println(err)
		return false
	}
	err=stmt.QueryRow(name).Scan(&id)
	if (err!=nil&&id!=upid) {
		fmt.Println(err)
		return false
	}
	return true
}
func upcomment(word string,id int)  {
	stmt,err:=DBconn().Prepare("update comment set message=? where id=?")
	if err!=nil {
		fmt.Println(err)
	}
	_,err=stmt.Exec(word,id)
	if err!=nil {
		fmt.Println(err)
	}
}
func deleter(id int)  {
	stmt,err:=DBconn().Prepare("delete from comment where id=?")
	if err!=nil {
		fmt.Println(err)
	}
	_,err=stmt.Exec(id)
	if err!=nil {
		fmt.Println(err)
	}
}


func Registe(c *gin.Context)  {
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	if UserSignup(username,password) {
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"注册成功"})
	}else {
		c.JSON(500,gin.H{"status":http.StatusInternalServerError,"massage":"数据库Insert报错"})
	}
}
func Login(c *gin.Context)  {
	username:=c.PostForm("username")
	password:=c.PostForm("password")
	if UserSignin(username,password) {
		c.SetCookie("username", username, 10, "/", "localhost", false, true)
		c.JSON(200,gin.H{"status":http.StatusOK,"message":"登录成功"})
	}else {
		c.JSON(403,gin.H{"status":http.StatusForbidden,"message":"登录失败，用户名或密码错误"})
	}
}
func SendMsg(c *gin.Context){
	if Cookie,err:=c.Request.Cookie("username");err==nil {
		username:=Cookie.Value
		if query(username) {
			message:=c.PostForm("message")
			pid:=c.PostForm("uid")
			uid,_:=strconv.Atoi(pid)
			if SendMesage(username,message,uid){
				c.JSON(200, gin.H{"内容":message,"用户名":username,"对id回复":uid})
			}else {
				c.JSON(403, gin.H{"status": http.StatusForbidden, "message": "发送失败"})
			}
		}
	}else {
		c.String(500,"cookie读取失败")
	}
}
func GetMsg(c *gin.Context){
	a:=c.Param("page")
	i, _ :=strconv.Atoi(a)
	floor:=GetMesage(i)
	for u:=0;u<10 ;u++  {
		for m,n:=range floor[u] {
			c.JSON(200,gin.H{"回复目标":u,"name":m,"message":n})
		}
	}
	if Cookie,err:=c.Request.Cookie("username");err==nil{
		name:=Cookie.Value
		if query(name) {
			fmt.Println("验证通过")
			if c.PostForm("like")=="1" {
				c.String(200,"你可以对心仪id点赞了")
				like(c.PostForm("id"))
			}
		}else {
			c.String(200,"验证失败")
			fmt.Println("name验证失败")
		}
	}else{
		fmt.Println("没有cookie")
	}
}
func Revoke(c *gin.Context){
	c.SetCookie("username", "0", -1, "/", "localhost", false, true)
	c.String(200,"你成功撤销登录")
}
func update(c *gin.Context)  {
	if Cookie,err:=c.Request.Cookie("username");err==nil {
		username:=Cookie.Value
		upid:=c.PostForm("id")
		id,_:=strconv.Atoi(upid)
		if affirm(id,username) {
			upword:=c.PostForm("word")
			upcomment(upword,id)
			c.String(200,"你更新了评论:"+upword )
		}else {
			c.String(200,"这不是你的评论 ")
		}
	}else {
		c.String(500,"cookie读取失败")
	}
}
func todelete(c *gin.Context)  {
	if Cookie,err:=c.Request.Cookie("username");err==nil {
		username:=Cookie.Value
		upid:=c.Param("id")
		id,_:=strconv.Atoi(upid)
		if affirm(id,username) {
			deleter(id)
			c.String(200,"你删除了你的评论")
		}else {
			c.String(200,"这不是你的评论 ")
		}
	}else {
		c.String(500,"cookie读取失败")
	}
}
func help(c *gin.Context)  {
	c.JSON(200,gin.H{"/registe":"注册账号使用（用户名不能相同）",
		"/login":"登录已有账号",
		"/sendmsg":"uid对目标id回复，message是发表内容",
		"/getmsg/:page":"id为点赞目标，like必须为1",
		"/revoke":"注销登录",
		"/update":"id为更新id，word为更新内容",
		"DELETE/msgs/:id":"id为删除目标"})
}

func main()  {
	r:=gin.Default()
	r.GET("/help",help)
	r.POST("/registe",Registe)
	r.POST("/login",Login)
	r.POST("/sendmsg",SendMsg)
	r.POST("/getmsg/:page",GetMsg)
	r.GET("/revoke",Revoke)
	r.PUT("/update",update)
	r.DELETE("DELETE/msgs/:id",todelete)
	r.Run()
}