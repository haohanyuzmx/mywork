package fundation

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
)
var DB *gorm.DB

func init()  {
	mysql,err:=gorm.Open("mysql","root:@(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err!=nil{
		log.Fatal(err,"connect sql")
	}
	DB=mysql
}
func Wrong(err error,mess string)  {
	if err!=nil {
		fmt.Println(err,mess)
	}
}