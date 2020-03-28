package main

import (
	"mygin/mygin"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
} 


func main()  {
	a:=mygin.Default()
	a.GET("/233", func(c *mygin.Context) {
		c.JSON(mygin.H{
			"233":233,
		})
	})
	a.Use(mid)
	a.POST("/login", func(c *mygin.Context) {
		/*var u user
		c.BindJSON(&u)
		c.JSON(u)*/
		a:=c.PostForm("123")
		c.String(a)
	})
	a.Run(":8080")
}
func mid(c *mygin.Context)  {
	a:=c.Getheader("Authorization")
	c.Set("user",a[7:])
}



/*func main()  {
	http.HandleFunc("/",value)
	http.ListenAndServe(":8080",nil)
}

func value(w http.ResponseWriter,r *http.Request)  {
	a:=r.Header["Authorization"]
	fmt.Fprintln(w, a)
}*/
