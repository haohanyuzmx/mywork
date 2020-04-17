package main

import (
	"Reptile-curriculum/timetable"
	"time"
)



func main()  {
	for i := 1; i <= 18; i++ {
		go timetable.Table(i)
	}
	for i1 := 2019210; i1 <=2019215 ; i1++ {
		for i2:=0;i2<=999;i2++ {
			i:=i1*1000+i2
			if i>2019215203 {
				goto OUT
			}
			go timetable.Stu(i)
		}
	}
OUT:	time.Sleep(15*time.Minute)
}
