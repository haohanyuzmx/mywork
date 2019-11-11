package main

import "fmt"

func main() {
	class1 := map[int]string{
		2003: "haha",
		2004: "hehe",
		2005: "baba",
	}
	class2 := map[int]string{
		2006: "guolai",
		2007: "guoqu",
	}
	tongxing := map[int]map[int]string{
		1: class1,
		2: class2,
	}
	fmt.Println(tongxing[1][2003])
}
