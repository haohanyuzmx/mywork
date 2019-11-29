package main

import (
	"database/sql"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)
import _ "github.com/go-sql-driver/mysql"
func main()  {
	enhen,_:=sql.Open("mysql","root:@tcp(localhost:3306)/assignment?charset=utf8")
	_,_=enhen.Prepare("create table students (id int(1),name varchar(10))")
	for i:=2019201303;i<=2019210346 ;i++  {
		a,_:=http.Get("http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php?xh="+strconv.Itoa(i))
		b,_:=ioutil.ReadAll(a.Body)
		re,_:=regexp.Compile("i....")
		c:=string(b)
		c=re.FindString(c)
		c=c[10:]
		_,_=enhen.Exec("insert into students values(i"+"'"+c+"')")
	}
}