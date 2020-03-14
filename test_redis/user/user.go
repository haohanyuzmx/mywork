package user

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	fundation "test_redis"
	jwt2 "test_redis/jwt"
)

type Userinfor struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Registe(c *gin.Context)  {
	var uf Userinfor
	if err:=c.ShouldBindJSON(&uf);err!=nil{
		c.JSON(200,gin.H{
			"code":500,
			"mess":"数据不全",
		})
		return
	}
	var u fundation.User
	fundation.DB.Where("username",uf.Username).First(&u)
	if u.ID!=0 {
		c.JSON(200,gin.H{
			"Code":500,
			"mess":"账号重合",
		})
		return
	}
	u.Username=uf.Username
	u.Password=uf.Password
	fundation.DB.Create(&u)
	c.JSON(200,gin.H{
		"code":200,
		"mess":"成功",
	})
}
func Login(c *gin.Context)  {
	var uf Userinfor
	if c.ShouldBindJSON(&uf)!=nil {
		c.JSON(200,gin.H{
			"code":500,
			"mess":"数据不全",
		})
		return
	}
	var u fundation.User
	fmt.Println(uf)
	err:=fundation.DB.Where(&fundation.User{Username:uf.Username,Password:uf.Password}).First(&u).Error
	if u.ID==0 {
		c.JSON(200,gin.H{
			"code":500,
			"mess":"错误",
		})
		fmt.Println(err)
		return
	}
	var jwt jwt2.JWT
	jwt.New(u)
	c.JSON(200,gin.H{
		"code":200,
		"token":jwt.Token,
	})
}
func Matchon(c *gin.Context)  {
	username:=c.Keys["username"]
	fmt.Println(username)
	conn:=fundation.Pool.Get()
	conn.Do("zadd","charts",0,username)
	c.JSON(200,gin.H{
		"code":200,
		"mess":"参与成功",
	})
}
func Matchoff(c *gin.Context )  {
	username:=c.Keys["username"]
	fmt.Println(username)
	conn:=fundation.Pool.Get()
	if _,err:=conn.Do("zrank",username);err!=nil {
		c.JSON(200,gin.H{
			"code":500,
			"mess":"你没有参加",
		})
		return
	}
	conn.Do("zrem","charts",username)
	c.JSON(200,gin.H{
		"code":200,
		"mess":"退出成功",
	})
}
func Charts(c *gin.Context)  {
	conn:=fundation.Pool.Get()
	resl,_:=redis.StringMap(conn.Do("zrevrange","charts",0,10,"withscores"))
	s:=make([]int ,len(resl))
	k:=make(map[int]string)
	ii:=0
	for i, i2 := range resl {
		value,err:=strconv.Atoi(i2)
		if err!=nil {
			fmt.Println(err)
		}
		s[ii]=value
		k[value]=i
		ii++
	}
	sort.Ints(s)
	var charts []gin.H
	for _, i2 := range s {
		chart:=gin.H{
			"usrename":k[i2],
			"scor":i2,
		}
		charts=append(charts,chart)
	}
	c.JSON(200,charts)
}
func Vote(c *gin.Context)  {
	username:=c.Keys["username"]
	voteuser:=c.PostForm("voteuser")
	conn:=fundation.Pool.Get()
	conn.Do("setnx",username,3)
	tm:=fundation.Overdue()
	conn.Do("expireat",username,tm)
	conn.Do("decr",username)
	conn.Do("incrby","charts",1,voteuser)
}
