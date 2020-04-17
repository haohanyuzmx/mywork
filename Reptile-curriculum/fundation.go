package fundation

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"regexp"
)

var DB *gorm.DB


func init()  {
	mysql,err:=gorm.Open("mysql","root:@(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	Wrong(err,"连接数据库")
	DB=mysql
}
func Wrong(err error,mess string)  {
	if err!=nil {
		fmt.Println(err,mess)
	}
}
func Getmess(now string,match string) ([][]string) {
	matchstring:=`<`+match+`.*?>`+`(.*?)`+`</`+match+`>`
	tomatch:=regexp.MustCompile(matchstring)
	return tomatch.FindAllStringSubmatch(now,-1)
}