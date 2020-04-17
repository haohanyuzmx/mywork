package timetable

import (
	fundation "Reptile-curriculum"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Class struct {
	Id string `gorm:"primary_key"`
	Types string
	Name string
	Teacher string
	Room string
	Time string
	School string
}

func Table(t int)  {
	urlm:="http://jwc.cqupt.edu.cn/kebiao/kb_stuYxh.php?type=stuYxNj&yxh="
	var st string
	if t<10 {
		st="0"+strconv.Itoa(t)
	}else {
		st=strconv.Itoa(t)
	}
	url:=urlm+st+"&nj=2019"
	hd,err:=http.Get(url)
	fundation.Wrong(err,"连接url")
	defer hd.Body.Close()
	all,err:=ioutil.ReadAll(hd.Body)
	fundation.Wrong(err,"body读取" )
	if hd.StatusCode!=200 {
		fmt.Println(hd.StatusCode,"连接url后出现问题")
	}
	allno:=strings.ReplaceAll(string(all),"\r\n","")
	alltr:=fundation.Getmess(allno,"tr")
	if len(alltr)==1 {
		return
	}
	c:=Class{Id: "",}
	for i1, i2 := range alltr {
		if i1>0 {
			alltd:=fundation.Getmess(i2[1],"td")
			if len(alltd)==10 {
				if c.Id!="" {
					fmt.Println(c)
					fundation.DB.Create(&c)
				}
				c.Types=alltd[0][1]
				c.Name=alltd[1][1]
				c.Id=alltd[2][1]
				c.Types+=alltd[3][1]
				c.School=alltd[4][1]
				c.Teacher=alltd[5][1]
				c.Time=alltd[6][1]
				c.Room=alltd[7][1]
			}else {
				c.Time+="和"+alltd[1][1]
				c.Room+="和"+alltd[2][1]
			}
		}
	}
}
