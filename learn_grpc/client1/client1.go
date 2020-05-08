package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	protol "learn_grpc/async"
	"time"
)

func main()  {
	conn,err:=grpc.Dial("localhost:55000",grpc.WithInsecure())
	if err!=nil {
		fmt.Println(err)
	}
	defer conn.Close()
	c:=protol.NewWebClient(conn)

	client,err:=c.Register(context.Background())
	if err!=nil {
		fmt.Println(err,"连接")
	}
	for i := 0; i < 20; i++ {
		go func() {
			u:=protol.User{
				Id:       int32(i),
				Name:     "433",
				Password: "12345",
			}
			err=client.Send(&u)
			if err!=nil {
				fmt.Println(err,"发送")
			}
			mess,err:=client.Recv()
			fmt.Println(mess,err)
		}()
	}
	time.Sleep(10*time.Second)
}

