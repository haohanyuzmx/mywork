package main

import (
	"github.com/gin-gonic/gin"
	"test_1/user"
)

func main() {
	app:=gin.Default()
	user.SetupRouter(app)
	app.Run()
}
