package main

import (
	"fmt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	pb "weifuwu/proto"
	"log"
	"net"
)

const (
	port = ":50051"
)


type server struct{}

func(s *server) Hello(ctx context.Context, in *pb.Request) (*pb.Response, error){
	fmt.Println("######## get client request name :" + in.Name)
	return &pb.Response{Greeting: "Hello " + in.Name}, nil
}


func main(){
	lis, err := net.Listen("tcp", port)
	if err != nil{
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterIrisFirstServer(s, &server{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil{
		log.Fatalf("failed to serve: %v", err)
	}
}

