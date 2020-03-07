package jwt

import (
	"fmt"
	"testing"
)

func TestChecktoken(t *testing.T) {
	a,b:=Checktoken(Creat("123"))
		if b==nil {
		fmt.Println(a)
	}else {
		fmt.Println(b)
	}
}
