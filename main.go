package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mywork/fundation"
	"mywork/web"
	"net/http"
	"strings"
)

func main() {
	fundation.Init()
	r:=gin.Default()
	r.Use(Cors())
	r.Static("/question","E:/GoProjects/src/mywork/templates/zhihu/question")
	r.Static("/comment","E:/GoProjects/src/mywork/templates/zhihu/comment")
	user:=r.Group("user")
	{
		user.POST("login",web.Login)
		user.POST("registe",web.Registe)
		user.GET("person",web.Mapa)
		user.POST("update",web.Update)
	}
	main:=r.Group("main")
	{
		main.GET("fir",web.Firpage)
		main.POST("fir",web.Select)
		main.GET("prom",web.Page)
		main.GET("like",web.Likequ)
		main.POST("mess",web.PutComment)
		main.POST("putqu",web.PutQuestion)
	}
	r.Run()
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		var headerKeys []string
		for k, _ := range c.Request.Header {
			headerKeys = append(headerKeys, k)
		}
		headerStr := strings.Join(headerKeys, ", ")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*")
		}
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "Options Request!")
		}
		c.Next()
	}
}