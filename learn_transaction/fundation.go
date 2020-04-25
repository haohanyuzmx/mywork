package fundation

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)

var DB *gorm.DB

func inti()  {
	mysql,err:=gorm.Open("mysql","root:@(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if Wrong(err,"connect db err") {
		return
	}
	DB=mysql
}
func Wrong(err error,mess string) (bool) {
	if err!=nil {
		log.Println(err,mess)
		return true
	}
	return false
}