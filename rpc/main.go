package main

import (
	hello "basic-server/proto"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
}

func main()  {
	g := grpc.NewServer()
	s := Server{}
	hello.RegisterGreeterServer(g,&s)
	lis, err := net.Listen("tcp", fmt.Sprintf(":8880"))
	if err != nil {
		panic("failed to listen: "+err.Error())
	}
	g.Serve(lis)

}

func (s *Server)  SayHello(ctx context.Context,request *hello.HelloRequest)(*hello.HelloReply,error){
	return &hello.HelloReply{Message:"Hello "+request.Name},nil
}

func (s *Server)  SayHelloAgain(ctx context.Context,request *hello.HelloRequest)(*hello.HelloReply,error){
	return &hello.HelloReply{Message:"Hello Again "+request.Name},nil
}