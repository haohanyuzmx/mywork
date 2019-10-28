package main
import "fmt"
func generate(ch chan int) {
	for i2 := 2; ; i2++ {
		ch <- i2
	}
}
func filter(in, out chan int, prime int) {
	for {
		i := <-in
		if i%prime != 0 {
			out <- i
		}
	}
}
func main() {
	ch := make(chan int)
	go generate(ch)
	for i1:=0;i1<100;i1++{
		prime := <-ch
		fmt.Print(prime, " ")
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}