package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
)
var m=0
func main() {
	m1 := make(map[string]int)
	m2:=make(map[int]string)
	var key []int
	for i:=2017;i<=2019 ;i++  {
		var m=209
	haha:		for ;m<=250 ;  {
		m++
		b=0
		for l:=0;l<=9 ;l++  {
			for h := 0; h <= 9; h++ {
				for j := 0; j <= 9; j++ {
					k := i*1000000 +1000*m+100*l +10*h + j
					a := strconv.Itoa(k)
					res, _ := http.Get("http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php?xh=" + a)
					defer res.Body.Close()
					body, _ := ioutil.ReadAll(res.Body)
					reg2 := regexp.MustCompile(a + "....")
					result2 := reg2.FindString(string(body))
					if result2[10:]=="  </"{
						b++
						if b==3 {
							goto haha
						}
					}else {
						b=0
						m1[result2] = 1
					}
				}
			}
		}
	}
	}
	fmt.Println(m1)
	for m,_:=range m1{
		for k,_:=range m1{
			if m[10:]==k[10:] {
				m1[m]++
			}
		}
	}
	for m,v:=range m1{
		m2[v]=m
	}
	for k:=range m2{
		key=append(key,k)
	}
	sort.Ints(key)
	for _,k:=range key{
		fmt.Println(k,m2[k])
	}
}


/*package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"sync"
	"time"
)

var lock sync.Mutex
var b=0
func main() {
	m2:=make(map[int]string)
	var key []int
	m1 := make(map[string]int)
	for i:=2017;i<=2019 ;i++  {
		var m  =209
		go shuchu(i,m,m1)
	}
	time.Sleep(5*time.Minute)
	caculate(m1)
	for m,v:=range m1{
		m2[v]=m
	}
	for k:=range m2{
		key=append(key,k)
	}
	sort.Ints(key)
	for _,k:=range key{
		fmt.Println(k,m2[k])
	}
}


func shuchu(i,m int,m1 map[string]int){
haha:		for ;m<=250 ;  {
	m++
	b=0
	for l:=0;l<=9 ;l++  {
		for h := 0; h <= 9; h++ {
			for j := 0; j <= 9; j++ {
				k := i*1000000 +1000*m+100*l +10*h + j
				a := strconv.Itoa(k)
				res, _ := http.Get("http://jwzx.cqupt.edu.cn/kebiao/kb_stu.php?xh=" + a)
				defer res.Body.Close()
				body, _ := ioutil.ReadAll(res.Body)
				reg2 := regexp.MustCompile(a + "....")
				result2 := reg2.FindString(string(body))
				if result2[10:]=="  </"{
					b++
					if b==3 {
						goto haha
					}
				}else {
					b=0
					lock.Lock()
					m1[result2] = 1
					lock.Unlock()
				}
			}
		}
	}
}
}

func caculate(m1 map[string]int){
	for m,_:=range m1{
		for k,_:=range m1{
			if m[10:]==k[10:] {
				m1[m]++
			}
		}
	}
}*/