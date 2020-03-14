package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	fundation "test_redis"
	"time"
)

type Header struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}
type Payload struct {
	Iss      string `json:"iss"`
	Exp      string `json:"exp"`
	Iat      string `json:"iat"`
	Username string `json:"username"`
}
type JWT struct {
	Header Header
	Payload Payload
	Signature string
	Token string
}

func (jwt *JWT) New(user fundation.User)  {
	jwt.Header=Header{
		Alg:"HS256",
		Typ:"JWT",
	}
	jwt.Payload=Payload{
		Iss:      "redrock",
		Exp:      strconv.FormatInt(time.Now().Add(3*time.Hour).Unix(), 10),
		Iat:      strconv.FormatInt(time.Now().Unix(), 10),
		Username: user.Username,
	}
	h,err:=json.Marshal(jwt.Header)
	if fundation.Wrong(err,"jwt.new"){
		return
	}
	p,err:=json.Marshal(jwt.Payload)
	if fundation.Wrong(err,"jwt.new"){
		return
	}
	baseh:=base64.StdEncoding.EncodeToString(h)
	basep:=base64.StdEncoding.EncodeToString(p)
	secret:=baseh+"."+basep
	key := "redrock"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(secret))
	s := mac.Sum(nil)
	jwt.Signature=base64.StdEncoding.EncodeToString(s)
	jwt.Token=secret+"."+jwt.Signature
}
func (jwt *JWT) Check (token string) (err error) {
	err = errors.New("token error")
	arr := strings.Split(token, ".")
	if len(arr)<3 {
		fundation.Wrong(err,"jwt.chek 1")
		return
	}
	baseh:=arr[0]
	h,err:=base64.StdEncoding.DecodeString(baseh)
	if fundation.Wrong(err,"jwt.chek 1"){
		return
	}
	json.Unmarshal(h,&jwt.Header)
	basep:=arr[1]
	p,err:=base64.StdEncoding.DecodeString(basep)
	if fundation.Wrong(err,"jwt.chek 1"){
		return
	}
	json.Unmarshal(p,&jwt.Payload)
	bases:=arr[2]
	s1,err:=base64.StdEncoding.DecodeString(bases)
	se:=baseh+"."+basep
	w:=[]byte("redrock")
	mac := hmac.New(sha256.New, w)
	mac.Write([]byte(se))
	s2 := mac.Sum(nil)
	if string(s1)!=string(s2) {
		fundation.Wrong(err,"token被改")
		fmt.Println("失败")
		return
	}else {
		fmt.Println("成功")
		jwt.Signature=arr[2]
		jwt.Token=token
	}
	return
}
