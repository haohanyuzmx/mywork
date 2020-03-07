package user

import (
	"fmt"
	"github.com/jinzhu/gorm"
	fundation "test_1"
	"testing"
)

func Test(t *testing.T) {
	//column:="password"
	//value:="555"
	var u =fundation.User{
		Model:    gorm.Model{
			ID:2,
		},
		Username: "wocao",
		Password: "666",
	}
	err:=fundation.DB.Model(&u).Update(fundation.User{Username:"n6"}).Error
	fmt.Println(err)
}
