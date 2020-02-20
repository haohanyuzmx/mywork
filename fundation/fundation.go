package fundation

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"mime/multipart"
	"os"
)
type Message struct {
	Id int
	Name string
	Comment string
	Path string
	ChildMessage *[]Message
}

var db *sql.DB

func Init()  {
	db,_=sql.Open("mysql","root:@tcp(localhost:3306)/zhihu?charset=utf8")
	db.SetMaxOpenConns(1000)
	err:=db.Ping()
	if err!=nil {
		fmt.Println("fail to connnect to db")
		fmt.Println(err.Error())
	}
}
func DBconn() *sql.DB {
	return db
}

func JsonNested( messageSlice []Message) []gin.H {
	var messageJsons []gin.H
	var messageJson gin.H
	for _, messages := range messageSlice {
		message := *messages.ChildMessage
		if messages.ChildMessage != nil {
			messageJson = gin.H{
				 "id" : messages.Id,
				"name":         messages.Name,
				"message":         messages.Comment,
				"picpath" :      messages.Path,
				"ChildrenMessage": JsonNested(message),
			}
		} else {
			messageJson = gin.H{
				"id" : messages.Id,
				"name": messages.Name,
				"message": messages.Comment,
				"picpath" :      messages.Path,
				"ChildrenMessage":"null",
			}
		}
		messageJsons = append(messageJsons, messageJson)
	}
	return messageJsons
}

func SaveUploadedFile(file *multipart.FileHeader,path string) error {
	dst := file.Filename
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	out, err := os.Create(path + dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, src)

	return err
}