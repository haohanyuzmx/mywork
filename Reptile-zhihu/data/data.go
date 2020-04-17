package data

import (
	fundation "Reptile-zhihu"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type data struct {
	Picture  string `json:"picture"`
	Name     string `json:"name"`
	Score    string `json:"score"`
	Director string `json:"director"`
	Comment  string `json:"comment"`
}

func TOP() ([]byte) {
	userAgent:="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/80.0.3987.163 Safari/537.36"
	url:="https://movie.douban.com/top250"
	c:=http.Client{
		Transport:     &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify:true,
			},
		},
	}
	var datas []data
	for i:=0;i<10 ;i++  {
		tail:="?start="+strconv.Itoa(i*25)+"&filter="
		urlr:=url+tail
		req,err:=http.NewRequest("GET",urlr,nil)
		fundation.Wrong(err,"make request")
		req.Header.Add("User-Agent",userAgent)
		resp,err:=c.Do(req)
		fundation.Wrong(err,"client wrong")
		defer resp.Body.Close()
		if resp.StatusCode!=200 {
			fmt.Println("not ok")
		}
		mess,err:=ioutil.ReadAll(resp.Body)
		if err!=nil {
			fmt.Println(err)
		}
		messno:=strings.ReplaceAll(string(mess),"\n","")
		ol:=regexp.MustCompile(`<ol (.*?)</ol>`)
		allist:=ol.FindAllStringSubmatch(messno,-1)
		li:=regexp.MustCompile(`<li>(.*?)</li>`)
		allthing:=li.FindAllStringSubmatch(allist[0][1],-1)
		for i1, i2 := range allthing {
			picandna:=regexp.MustCompile(`<img width="100" alt="(.*?)" src="(.*?)" class="">`)
			pn:=picandna.FindStringSubmatch(i2[1])
			dir:=regexp.MustCompile(`<p class="">                            导演:(.*?)&nbsp`)
			com:=regexp.MustCompile(`<span class="inq">(.*?)</span>`)
			soc:=regexp.MustCompile(`<span class="rating_num" property="v:average">(.*?)</span>`)
			s:=soc.FindStringSubmatch(i2[1])
			var c string
			if i==9 {
				if i1==19||i1==22||i1==24 {
					c=""
				}
			}else {
				c=com.FindStringSubmatch(i2[1])[1]
			}
			d:=dir.FindStringSubmatch(i2[1])
			data:=data{
				Picture:  pn[2],
				Name:     pn[1],
				Director: d[1],
				Comment:  c, //244,247,249没有评语
				Score:    s[1],
			}
			datas=append(datas,data)
		}
	}
	a,err:=json.Marshal(datas)
	fundation.Wrong(err,"json格式化")
	return a
}