package main

import ("fmt"
	"time"
)
var ch =make(chan int)
var i int
func caculate() (){
	i = <-ch
	for ; ; i++ {
		ok:=0
		for c := 2; c < i; c++ {
			if i%c == 0 {
				ok=1
				break
			}

		}
		if ok==0 {
			fmt.Println(i)
			ch<-i
		}
	}
}
func main() {
	go func() {ch<-2}()
	for   {
		go caculate()
		time.Sleep(1)
		if i>1000 {
			break
		}
	}
}
