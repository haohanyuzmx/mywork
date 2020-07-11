package fundation

import (
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)
type User struct {
	gorm.Model
	Username string
	Password string
}



var Pool redis.Pool
var DB *gorm.DB
func init()  {
	mysql,err:=gorm.Open("mysql","root:@(127.0.0.1:3306)/ assessment?charset=utf8&parseTime=true")
	if err!=nil {
		log.Println("mysql连接失败")
		return
	}
	DB=mysql
	Pool = redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 120,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
	if Pool.Get()==nil {
		log.Println("redis连接失败")
		return
	}
}
func Wrong(err error,where string) bool {
	if err!=nil {
		log.Println(where)
		log.Println(err)
		return true
	}
	return false
}