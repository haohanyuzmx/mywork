package user

import (
	"github.com/gin-gonic/gin"
	fundation "learn_transaction"
	"regexp"
	"strconv"
	"strings"
)

type Student struct {
	Id string `gorm:"primary_key"`
	Time string/*ABCDE表示周数,123456表示第几节,(n1-n2(n3))代表起始时间,n3=1是单周,等于2是双，0是全部,+号连接
	例如A5(2-16（2))就是周一五节课从第2周到16周的双周*/
	Classid string `gorm:"index:class_id"`
}
type Class struct {
	Id string `gorm:"primary_key"`
	Types string
	Name string
	Teacher string
	Room string
	Time string
	School string
	Num int
	Max int
}


func Chose(c *gin.Context)  {
	var status int
	id:=c.PostForm("id")
	classid:=c.PostFormArray("classid")
	tx:=fundation.DB.Begin()
	stu:=Student{Id: id}
	tx.Find(&stu)
	stutime :=stu.Time
	stuclassid:=stu.Classid
	timeinfor :=regexp.MustCompile(`(.*?)\((.*?)\((.*?)\){2}`)
	for _, i2 := range classid {
		allstutime :=strings.Split(stutime,"+")
		status=0
		cls:=Class{Id: i2}
		tx.Find(&cls)
		clatime:=cls.Time
		if cls.Num<cls.Max {
			allclatime:=strings.Split(clatime,"+")
			for _, i3 := range allstutime {
				stuinfor:= timeinfor.FindStringSubmatch(i3)
				for _, i5 := range allclatime {
					clainfor:=timeinfor.FindStringSubmatch(i5)
					if stuinfor[1]==clainfor[1] {
						week1:=strings.Split(stuinfor[2],"-")
						week1i,_:=strconv.Atoi(week1[1])
						week2:=strings.Split(clainfor[2],"-")
						week2i,_:=strconv.Atoi(week2[0])
						if week1i>week2i {
							if stuinfor[3]=="0"||clainfor[3]=="0"||stuinfor[3]==clainfor[3] {
								status=1
								c.JSON(200,gin.H{
									"课程重叠":cls.Id,
								})
							}
						}
					}
				}
			}
		}else {
			status=1
			c.JSON(200,gin.H{
				"课程已满":cls.Id,
			})
		}
		if status==0 {
			stutime+=clatime
			stuclassid+=cls.Id
			cls.Num++
			tx.Update(&cls)
		}
	}
	stu.Classid=stuclassid
	stu.Time=stutime
	tx.Update(&stu)
	tx.Commit()
}