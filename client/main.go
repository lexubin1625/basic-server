package main

import (
hello "basic-server/proto"
"context"
"fmt"
"google.golang.org/grpc"
)

func main()  {
	conn,err := grpc.Dial("127.0.0.1:8880",grpc.WithInsecure())
	if err!=nil{
		panic(err)
	}
	defer conn.Close()
	c := hello.NewGreeterClient(conn)
	r,err := c.SayHello(context.Background(),&hello.HelloRequest{Name:"ucan"})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r.Message)
	r,err = c.SayHelloAgain(context.Background(),&hello.HelloRequest{Name:"ucan"})
	if err!=nil{
		panic(err)
	}
	fmt.Println(r.Message)
}
