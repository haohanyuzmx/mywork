package mygin

import (
	"log"
	"net/http"
	"strings"
)

type Hander func(*Context)
type handermap map[string][]Hander
type App struct {
	router map[string]handermap
}

func Default() *App {
	return &App{router: make(map[string]handermap)}
}

func (a *App)Use(hander Hander)  {
	for _, i2 := range a.router {
		for _, i4 := range i2 {
			i4=append(i4,hander)
		}
	}
}

func (a *App)GET(url string,hander Hander)  {
	a.hander("GET",url,hander)
}
func (a *App)POST(url string,hander Hander)  {
	a.hander("POST",url,hander)
}
func (a *App)Run(port string)  {
	http.Handle("/",a)
	if er:=http.ListenAndServe(port,nil);er!=nil {
		log.Fatal(er)
	}
}

func (a *App)hander(meth,url string,hander Hander)  {
	handers,ok:=a.router[meth]
	if !ok {
		m:=make(handermap)
		a.router[meth]=m
		handers=m
	}
	if  _,ok=handers[url];ok{
		panic("same url")
	}
	var h []Hander
	if h, ok = handers[url];!ok {
		h=make([]Hander,0)
	}
	h=append(h,hander)
	handers[url]=h
}



func (a *App)ServeHTTP(writ http.ResponseWriter,req *http.Request)  {
	httpmeth:=req.Method
	uri:=req.RequestURI
	uris:=strings.Split(uri,"?")
	if len(uris)<1 {
		return
	}

	handers,ok:=a.router[httpmeth]
	if !ok {
		log.Println("hacker")
		return
	}
	m,ok:=handers[uris[0]]
	if !ok {
		Handler404(writ,req)
		return
	}
	h:=NewContext(req,writ)
	for _, i2 := range m {
		if i2!=nil {
			i2(&h)
		}
	}
}

func Handler404(w http.ResponseWriter,req *http.Request)  {
	log.Println(404,req.RequestURI)
	w.Write([]byte("404 not find"))
}