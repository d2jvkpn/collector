package main

import (
	// "context"
	"flag"
	"fmt"
	"log"

	"github.com/d2jvkpn/collector/proto"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn *grpc.ClientConn
	cli  proto.RecordServiceClient
}

func NewGrpcClient(conn *grpc.ClientConn) *GrpcClient {
	return &GrpcClient{
		conn: conn,
		cli:  proto.NewRecordServiceClient(conn),
	}
}

func main() {
	var (
		addr string
		err  error
		// ctx    context.Context
		conn   *grpc.ClientConn
		client *GrpcClient
	)

	flag.StringVar(&addr, "addr", "grpc address", "localhost:5021")
	flag.Parse()

	inte := proto.ClientInterceptor{
		Headers: map[string]string{},
	}

	conn, err = grpc.Dial(addr,
		grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(inte.Unary()),
		grpc.WithStreamInterceptor(inte.Stream()),
	)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(">> Dial:", addr)

	defer func() {
		conn.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	client = NewGrpcClient(conn)
	fmt.Println("~~~", client)
	// TODO:...
}
