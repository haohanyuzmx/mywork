package fundation

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)
var Pool redis.Pool
var DB *gorm.DB
var Tm int64=1583971200


type User struct {
	gorm.Model
	Username string
	Password string
}


func init()  {
	Pool = redis.Pool{
		MaxIdle:     16,
		MaxActive:   32,
		IdleTimeout: 120,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		},
	}
	mysql,err:=gorm.Open("mysql","root:@(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err!=nil {
		fmt.Println(err)
	}
	DB=mysql
}
func Wrong(err error,where string) bool {
	if err!=nil {
		fmt.Println(where)
		fmt.Println(err)
		return true
	}
	return false
}
func Overdue() int64 {
	for i:=0; ;i++  {
		if Tm<time.Now().Unix() {
			Tm+=86400
		}else {
			break
		}
	}
	return Tm
}
