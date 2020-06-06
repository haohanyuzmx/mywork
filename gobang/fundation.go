package fundation

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
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
		fmt.Println(err)
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
}
func Wrong(err error,where string) bool {
	if err!=nil {
		fmt.Println(where)
		fmt.Println(err)
		return true
	}
	return false
}