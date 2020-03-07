package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	fundation "test_1"
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


func NewHeader() Header {
	return Header{
		Alg: "HS256",
		Typ: "JWT",
	}
}
func Creat(username string) (token string) {
	header := NewHeader()
	payload := Payload{
		Iss:      "redrock",
		Exp:      strconv.FormatInt(time.Now().Add(3*time.Hour).Unix(), 10),
		Iat:      strconv.FormatInt(time.Now().Unix(), 10),
		Username: username,
	}
	h,err:=json.Marshal(header)
	fundation.Wrong(err)
	p,err:=json.Marshal(payload)
	fundation.Wrong(err)
	baseh:=base64.StdEncoding.EncodeToString(h)
	basep:=base64.StdEncoding.EncodeToString(p)
	secret:=baseh+"."+basep
	key := "redrock"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(secret))
	s := mac.Sum(nil)
	signature:=base64.StdEncoding.EncodeToString(s)
	token=secret+"."+signature
	return
}

func Checktoken(token string) (username string,err error) {
	err = errors.New("token error")
	arr := strings.Split(token, ".")
	if len(arr)<3 {
		return 
	}
	_,err=base64.StdEncoding.DecodeString(arr[0])
	if err!=nil {
		return 
	}
	pay,err:=base64.StdEncoding.DecodeString(arr[1])
	if err!=nil {
		return
	}
	s1,err:=base64.StdEncoding.DecodeString(arr[2])
	str:=arr[0]+"."+arr[1]
	key := "redrock"
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(str))
	s2 := mac.Sum(nil)
	if string(s1)!=string(s2) {
		return
	}
	var payload Payload
	json.Unmarshal(pay,&payload)
	username=payload.Username
	err=nil
	return
}