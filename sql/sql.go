package sql

import (
	"fmt"
	"mywork/fundation"
	"strconv"
)

func UserSignup(username ,passward,telephone string) (bool) {
	num,err:=strconv.Atoi(telephone)
	stmt,err:=fundation.DBconn().Prepare("insert into user(name,password,telephone) value (?,?,?)")
	if err!=nil{
		fmt.Println("fail to insert1")
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()
	_,err=stmt.Exec(username,passward,num)
	if err!=nil{
		fmt.Println("fail to insert2")
		fmt.Println(err.Error())
		return false
	}
	return true
}
func UserSignin(username string,passward string) (bool,string) {
	var id int
	var name string
	if username[:4]!="tele" {
		stmt,err:=fundation.DBconn().Prepare("select id from user where name=? and password=?")
		if err!=nil{
			fmt.Println(err.Error())
			return false,"0"
		}
		err=stmt.QueryRow(username,passward).Scan(&id)
		if err!=nil{
			fmt.Println(err.Error())
			return false,"0"
		}
		return true,"0"
	}
	stmt,err:=fundation.DBconn().Prepare("select name from user where telephone=? and password=?")
	if err!=nil{
		fmt.Println(err.Error())
		return false,"0"
	}
	telephone,err:=strconv.Atoi(username[4:])
	err=stmt.QueryRow(telephone,passward).Scan(&name)
	if err!=nil{
		fmt.Println(err.Error())
		return false,"0"
	}
	return true,name
}
func FindMessageByPid(pid int,quid string) []fundation.Message {
	db := fundation.DBconn()
	res, err := db.Query("select id,comment,name,picture from "+quid+" where pid=?", pid)
	if err != nil {
		fmt.Println("du.query findMessage is error !",err)
	}
	var messageSlice []fundation.Message
	for res.Next() {
		var messages fundation.Message
		err := res.Scan(&messages.Id, &messages.Comment, &messages.Name,&messages.Path)
		if err != nil {
			fmt.Println("res.scan id is error !", err)
		}
		child := FindMessageByPid(messages.Id,quid)
		messages.ChildMessage = &child
		messageSlice = append(messageSlice, messages)
	}
	return messageSlice
}
func Person(name,mold1,mold2 string) []string {
	db:=fundation.DBconn()
	stat,err:=db.Prepare("select "+mold1+" from perpage where "+mold2+"=?")
	wrong(err)
	res,err:=stat.Query(name)
	wrong(err)
	var persons []string
	for res.Next()  {
		var person string
		err:=res.Scan(&person)
		fmt.Println(person)
		wrong(err)
		persons=append(persons,person)
	}
	return persons
}
func Perportrait(name string) (portrait string) {
	db:=fundation.DBconn()
	stat,err:=db.Prepare("select portraint from user where name=?")
	wrong(err)
	res:=stat.QueryRow(name)
	err=res.Scan(&portrait)
	return
}
func Ques(name string) []string {
	db:=fundation.DBconn()
	stat,err:=db.Prepare("select quname from question where curiouser=?")
	wrong(err)
	res,err:=stat.Query(name)
	wrong(err)
	var qunames []string
	for res.Next()  {
		var quname string
		err:=res.Scan(&quname)
		wrong(err)
		qunames=append(qunames,quname)
	}
	return qunames
}
func Upinform(table,column,value,name string) bool {
	db:=fundation.DBconn()
	stat,err:=db.Prepare("update "+table+" set"+column+"=? where name=?")
	a:=wrong(err)
	if a {
		_,err=stat.Exec(value,name)
		a=wrong(err)
		if a {
			return true
		}
		return false
	}
	return false
}
func Firpagesort() ([]string,[]string,[]string,[]string,[]string) {
	db:=fundation.DBconn()
	stat,err:=db.Prepare("select quname,question,picture,intkey,likenum from question order by likenum desc")
	wrong(err)
	res,err:=stat.Query()
	wrong(err)
	var qunames,questions,pictures,keys,likenums []string
	for res.Next()  {
		var quname,question,picture,key,likenum string
		var intkey,intlike int
		err:=res.Scan(&quname,&question,&picture,&intkey,&intlike)
		wrong(err)
		key=strconv.Itoa(intkey)
		likenum=strconv.Itoa(intlike)
		likenums=append(likenums,likenum)
		keys=append(keys,key)
		qunames=append(qunames,quname)
		questions=append(questions,question)
		pictures=append(pictures,picture)
	}
	return qunames,questions,pictures,keys,likenums
}
func Firpageunsort() ([]string,[]string,[]string,[]string,[]string) {
	db:=fundation.DBconn()
	stat,err:=db.Prepare("select quname,question,picture,intkey,likenum from question order by rand()")
	wrong(err)
	res,err:=stat.Query()
	wrong(err)
	var qunames,questions,pictures,keys,likenums []string
	for res.Next()  {
		var quname,question,picture,key,likenum string
		var intkey,intlike int
		err:=res.Scan(&quname,&question,&picture,&intkey,&intlike)
		wrong(err)
		key=strconv.Itoa(intkey)
		likenum=strconv.Itoa(intlike)
		likenums=append(likenums,likenum)
		keys=append(keys,key)
		qunames=append(qunames,quname)
		questions=append(questions,question)
		pictures=append(pictures,picture)
	}
	return qunames,questions,pictures,keys,likenums
}
func PutMessage(id,name,comment,pid string) bool {
	db:=fundation.DBconn()
	intpid,err:=strconv.Atoi(pid)
	wrong(err)
	table:="comment_"+id
	stat,err:=db.Prepare("insert into "+table+" (name,comment,pid)value (?,?,?)")
	if wrong(err)==false{
		return false
	}
	_,err=stat.Exec(name,comment,intpid)
	if wrong(err)==false{
		return false
	}
	return true
}
func Likequ (id int)  {
	var likenum int
	db:=fundation.DBconn()
	stat,err:=db.Prepare("select likenum from question where id=?")
	wrong(err)
	err=stat.QueryRow(id).Scan(&likenum)
	wrong(err)
	likenum++
	stat,err=db.Prepare("update question set likenum=? where id=?")
	_,err=stat.Exec(likenum,id)
	wrong(err)
}
func Select(mess string) ([]string,[]string,[]string,[]string,[]string) {
	db:=fundation.DBconn()
	mess="%"+mess+"%"
	res,err:=db.Query("select quname,question,picture,intkey,likenum from question where quname like ?",mess)
	wrong(err)
	var qunames,questions,pictures,keys,likenums []string
	for res.Next()  {
		var quname,question,picture,key,likenum string
		var intkey,intlike int
		err:=res.Scan(&quname,&question,&picture,&intkey,&intlike)
		wrong(err)
		key=strconv.Itoa(intkey)
		likenum=strconv.Itoa(intlike)
		likenums=append(likenums,likenum)
		keys=append(keys,key)
		qunames=append(qunames,quname)
		questions=append(questions,question)
		pictures=append(pictures,picture)
	}
	return qunames,questions,pictures,keys,likenums
}
func Prom(key int) (string,string,string,string) {
	db:=fundation.DBconn()
	res:=db.QueryRow("select quname,question,picture,intkey,likenum from question where intkey=?",key)
	var quname,question,picture,likenum string
	var intkey,intlike int
	err:=res.Scan(&quname,&question,&picture,&intkey,&intlike)
	wrong(err)
	strconv.Itoa(intkey)
	likenum=strconv.Itoa(intlike)
	return quname,question,picture,likenum
}
func PutQuestion(quname,name,question string) bool {
	db:=fundation.DBconn()
	res,err:=db.Prepare("insert into question (quname,question,curiouser) value(?,?,?)")
	if wrong(err) {
		_,err=res.Exec(quname,question,name)
		if wrong(err){
			return true
		}
		return false
	}
	return false
}


func wrong(err error) bool {
	if err!=nil {
		fmt.Println(err)
		return false
	}
	return true
}