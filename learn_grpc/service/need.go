package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	protol "learn_grpc/async"
	"strconv"
	"strings"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
type Payload struct {
	Iss string `json:"iss"`
	Exp string `json:"exp"`
	Iat string `json:"iat"`
	Id  string `json:"id"`
}

var DB *gorm.DB

func init() {
	mysql, err := gorm.Open("mysql", "root:@(127.0.0.1:3306)/test?charset=utf8&parseTime=true")
	if err != nil {
		fmt.Println(err)
	}
	DB = mysql
}
func NewJwt(u protol.User) string {
	hea := Header{
		Alg: "HS256",
		Typ: "JWT",
	}
	pay := Payload{
		Iss: "redrock",
		Exp: strconv.FormatInt(time.Now().Add(3*time.Hour).Unix(), 10),
		Iat: strconv.FormatInt(time.Now().Unix(), 10),
		Id:  strconv.Itoa(int(u.GetId())),
	}
	h, err := json.Marshal(hea)
	Wrong(err, "json化header出错")
	p, err := json.Marshal(pay)
	Wrong(err, "json化pay错误")
	baseh := base64.StdEncoding.EncodeToString(h)
	basep := base64.StdEncoding.EncodeToString(p)
	secret := baseh + "." + basep
	key := "redrock"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(secret))
	s := mac.Sum(nil)
	token:=secret+"."+base64.StdEncoding.EncodeToString(s)
	return token
}

func Check(user *protol.User) (err error) {
	var myh Header
	var myp Payload
	err = errors.New("token error")
	arr := strings.Split(user.Login.Token, ".")
	if len(arr) < 3 {
		Wrong(err, "chek 1")
		return
	}
	baseh := arr[0]
	h, err := base64.StdEncoding.DecodeString(baseh)
	if  Wrong(err, "chek 1") {
		return
	}
	json.Unmarshal(h, &myh)
	basep := arr[1]
	p, err := base64.StdEncoding.DecodeString(basep)
	if  Wrong(err, "chek 2") {
		return
	}
	json.Unmarshal(p, &myp)
	bases := arr[2]
	s1, err := base64.StdEncoding.DecodeString(bases)
	se := baseh + "." + basep
	w := []byte("redrock")
	mac := hmac.New(sha256.New, w)
	mac.Write([]byte(se))
	s2 := mac.Sum(nil)
	if string(s1) != string(s2) {
		Wrong(err, "token被改")
		fmt.Println("失败")
		return
	} else {
		fmt.Println("成功")
		i,_:=strconv.Atoi(myp.Id)
		user.Id=int32(i)
	}
	return
}
func Wrong(err error, mess string) bool {
	if err != nil {
		fmt.Println(err, mess)
		return true
	}
	return false
}