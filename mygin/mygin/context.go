package mygin

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Context struct {
	req *http.Request
	write http.ResponseWriter
	queryParam map[string]string
	formParam map[string]string
	json []byte
	Key map[string]interface{}
}
type H map[string]interface{}

func NewContext(rep *http.Request,write http.ResponseWriter) (Context) {
	log.Println(rep.RequestURI)
	var c=Context{
		req:        rep,
		write:       write,
	}
	c.queryParam=parseQuery(rep)
	c.formParam=parseForm(rep)
	c.json=parseJson(rep)
	return c
}


func (c *Context)Getheader(key string) (string) {
	return c.req.Header[key][0]
}

func (c *Context) Set(key string,value interface{}) {
	if c.Key == nil {
		c.Key = make(map[string]interface{})
	}
	c.Key[key]=value
}

func (c *Context)Query(key string)(string){
	a:=c.queryParam[key]
	return a
}

func (c *Context)PostForm(key string) (string) {
	a:=c.formParam[key]
	return a
}

func (c *Context)BindJSON(s interface{}) error {
	return json.Unmarshal(c.json,s)
}




func (c *Context) String(s string)  {
	_,err:=c.write.Write([]byte(s))
	if err!=nil {
		log.Println(err)
	}
}

func (c *Context)JSON(obj interface{}) {
	a,err:=json.Marshal(obj)
	if err!=nil {
		log.Println(err)
	}
	c.write.Write(a)
}






func parseQuery(r *http.Request)(param map[string]string)  {
	param=make(map[string]string)
	uri:=r.RequestURI
	uris:=strings.Split(uri,"?")
	if (len(uris)==1) {
		return
	}
	thing:=strings.Split(uris[len(uris)-1],"&")
	for _, i2 := range thing {
		num:=strings.Split(i2,"=")
		if len(num)!=2 {
			return
		}
		param[num[0]]=num[1]
	}
	return
}

func parseForm(r *http.Request) (form map[string]string) {
	err:=r.ParseForm()
	if err!=nil {
		log.Println(err)
		return
	}
	for i, i2 := range r.PostForm {
		form[i]=i2[0]
	}
	return
}

func parseJson(r *http.Request) (a []byte) {
	a=make([]byte,1)
	a,err:=ioutil.ReadAll(r.Body)
	if err!=nil {
		log.Println(err)
		return
	}
	return
}