package main

import (
	"fmt"
	"time"
)

func factorial(n int, ch chan int) {
	var res = 1
	for i := 1; i <= n; i++ {
		res *= i
	}
	ch <- res
}
func main() {
	ch := make(chan int)
	go func() {
		for i := 1; i <= 20; i++ {
			factorial(i, ch)
		}
	}()
	go func() {
		for i := range ch {
			fmt.Println(i)
		}
	}()
	time.Sleep(2e9)
}
