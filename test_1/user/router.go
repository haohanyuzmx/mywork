package user

import (
	"github.com/gin-gonic/gin"
	"test_1/middleware"
)

func SetupRouter(app *gin.Engine)  {
	app.POST("login",Login)
	app.POST("register",Register)
	app.GET("all",All)
	app.POST("update",middleware.Inspect,Update)
}
