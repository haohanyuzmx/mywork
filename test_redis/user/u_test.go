package user

import (
	"fmt"
	fundation "test_redis"
	"testing"
)

func TestCharts(t *testing.T) {
	var u fundation.User
	uf:=Userinfor{
		Username: "123",
		Password: "233",
	}
	err:=fundation.DB.Where(&fundation.User{Username:uf.Username,Password:uf.Password}).First(&u).Error
	fmt.Println(err)
	fmt.Println(u)
}
