package main

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	protol "learn_grpc/async"
	"net"
)

type service struct {
	protol.UnimplementedWebServer
}

func (s *service)Register(conn protol.Web_RegisterServer) error  {
	for  {
		in,err:=conn.Recv()
		fmt.Println(in)
		if Wrong(err,"接收数据") {
			return err
		}
		var quer protol.User
		retu:=new(protol.Ret)
		DB.Where(&protol.User{Id: in.Id}).First(&quer)
		if quer.GetId()==in.GetId() {
			err:=errors.New("id重复")
			return err
		}
		DB.Create(in)
		retu.Message="成功"
		err=conn.Send(retu)
		if Wrong(err,"发送数据") {
			return err
		}
	}
	return nil
}
func (s *service)Login(conn protol.Web_LoginServer) error  {
	for  {
		in,err:=conn.Recv()
		if Wrong(err,"接收数据") {
			return err
		}
		err=DB.Where(&protol.User{Id: in.Id,Password: in.Password}).Error
		if Wrong(err,"token序列化错误") {
			return  err
		}
		jwt:=new(protol.Jwt)
		jwt.Token=NewJwt(*in)
		err=conn.Send(jwt)
		if Wrong(err,"发送数据") {
			return err
		}
	}
	return nil
}
func (s *service)Update(conn protol.Web_UpdateServer) error  {
	for  {
		in,err:=conn.Recv()
		if Wrong(err,"接收数据") {
			return err
		}
		err=Check(in)
		if err!=nil {
			fmt.Println(err)
			return  err
		}
		in.Login=nil
		err=DB.Model(in).Update(protol.User{
			Id:       in.Id,
			Name:     in.Name,
			Password: in.Password,
		}).Error
	}
}

func main(){
	lis,err:=net.Listen("tcp",":55000")
	if err!=nil {
		fmt.Println(err)
	}
	s:=grpc.NewServer()
	protol.RegisterWebServer(s,&service{})
	s.Serve(lis)
}