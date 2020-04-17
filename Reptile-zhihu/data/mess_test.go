package data

import (
	fundation "Reptile-zhihu"
	"fmt"
	"testing"
)

func TestTOP(t *testing.T) {
	fmt.Println(string(TOP()))
}
func TestTOP2(t *testing.T) {
	var data data
	err:=fundation.DB.CreateTable(&data).Error
	fmt.Println(err)
}