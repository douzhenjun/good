package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "weifuwu/proto"
	"log"
	"os"
)

const(
	address = "localhost:50051"
	defaultName = "world"
)


func main(){
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil{
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewIrisFirstClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	r, err := c.Hello(context.Background(), &pb.Request{Name: name})
	if err != nil{
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Greeting)
}