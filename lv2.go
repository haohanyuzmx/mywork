package main

import (
	"fmt"
	"strconv"
	"time"
)

func main() {
	timeLayout := "2006-01-02 15:04:05"
	var slice []int64
HAHA:
	var timestamp string
	fmt.Scanf("%s", &timestamp)
	if timestamp == "result" {
		goto HEHE
	} else {
		c, _ := strconv.ParseInt(timestamp, 10, 64)
		slice = append(slice, c)
		goto HAHA
	}
HEHE:
	for _, v := range slice {
		datetime := time.Unix(v, 0).Format(timeLayout)
		fmt.Println(datetime)
	}
}
