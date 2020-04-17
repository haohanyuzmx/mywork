package timetable

import (
	fundation "Reptile-curriculum"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type Student struct {
	Id string
	Name string
	Classid string `gorm:"primary_key"`
}

func Stu(id int)  {
	urlm:="http://jwc.cqupt.edu.cn/kebiao/kb_stu.php?xh="
	myid:=strconv.Itoa(id)
	s:=Student{Id:myid}
	url:=urlm+myid
	hd,err:=http.Get(url)
	fundation.Wrong(err,"连接url")
	defer hd.Body.Close()
	all,err:=ioutil.ReadAll(hd.Body)
	fundation.Wrong(err,"body读取" )
	if hd.StatusCode!=200 {
		fmt.Println(hd.StatusCode,"连接url后出现问题",id)
		return
	}
	allno:=strings.ReplaceAll(string(all),"\r\n","")
	na:=regexp.MustCompile(myid+"(.*?)</li>")
	s.Name=na.FindStringSubmatch(allno)[1]
	ytall:=fundation.Getmess(allno,"div")
	if len(ytall)<27 {
		return
	}
	yt:=fundation.Getmess(ytall[26][1],"tr")
	for i, i2 := range yt {
		if i!=0 {
			ytd:=fundation.Getmess(i2[1],"td")
			if len(ytd)==10 {
				s.Classid=ytd[2][1]
				fundation.DB.Create(&s)
				fmt.Println(s)
			}
		}
	}
}
