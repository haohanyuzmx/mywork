package user

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	fundation "gobang"
	"strconv"
	"strings"
	"testing"
)

func TestAddRoom(t *testing.T) {
	conn:=fundation.Pool.Get()
	//m,err:=redis.Strings(conn.Do("smembers","room1"))
	allstep, _ :=redis.StringMap(conn.Do("ZRANGEBYSCORE","room1+chess","-inf","+inf","withscores"))
	var thechess [16][16]int
	for i, i2 := range allstep {
		xy := strings.Split(i, ",")
		x, _ := strconv.Atoi(xy[0])
		y, _ := strconv.Atoi(xy[1])
		thechess[x][y], _ = strconv.Atoi(i2)
	}
	win := judge(thechess)
	if win != 0 {
		fmt.Println("1")
	}
	mjson, _ := json.Marshal(allstep)
	a := string(mjson)
	fmt.Println(allstep,"\n",mjson,"\n",a)
	//a:=strings.Split("1+11,12","+")
	//fmt.Println(a[0],a[1])
}
func TestChess(t *testing.T) {
	fundation.DB.CreateTable(&fundation.User{})
}
