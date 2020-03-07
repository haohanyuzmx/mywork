package fundation

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)
var DB *gorm.DB

type User struct {
	gorm.Model
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}


func init()  {
	mysql,err:=gorm.Open("mysql","root:@(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err!=nil {
		fmt.Println(err)
	}
	DB=mysql
}


func Wrong(err error) bool {
	if err!=nil {
		fmt.Println(err)
		return false
	}
	return true
}