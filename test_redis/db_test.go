package fundation

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"sort"
	"strconv"
	"testing"
)

func TestOverdue(t *testing.T) {
	red:=Pool.Get()
	_,err:=red.Do("zadd","charts",3,"zhansan")
	fmt.Println(err)
	_,err=red.Do("zadd","charts",0,"zaosi")
	_,err=red.Do("zadd","charts",20,"wangwu")
	resl,err:=redis.StringMap(red.Do("zrevrange","charts",0,10,"withscores"))
	fmt.Println(err)
	fmt.Println(resl)
	s:=make([]int ,len(resl))
	k:=make(map[int]string)
	ii:=0
	for i, i2 := range resl {
		value,err:=strconv.Atoi(i2)
		if err!=nil {
			fmt.Println(err)
		}
		s[ii]=value
		k[value]=i
		ii++
	}
	sort.Ints(s)
	var charts []gin.H
	for _, i2 := range s {
		chart:=gin.H{
			"usrename":k[i2],
			"scor":i2,
		}
		charts=append(charts,chart)
	}
	fmt.Println(charts)
}
func Test(t *testing.T)  {
	var a interface{}
	a=123
	h:=a.(string)
	fmt.Println(h)
	fmt.Println(a)
}

func TestOverdue2(t *testing.T) {
	println(Overdue())
}

