package user

import (
	"fmt"
	"regexp"
	"testing"
)

func Test(t *testing.T) {
	tim:="A5(2-16(2))"
	timeinfor:=regexp.MustCompile(`(.*?)\((.*?)\((.*?)\){2}`)
	infor:=timeinfor.FindStringSubmatch(tim)
	fmt.Println(infor)
}
