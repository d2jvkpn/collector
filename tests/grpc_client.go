package main

import (
	"context"
	"flag"
	// "fmt"
	"log"

	"github.com/d2jvkpn/collector/proto"

	"google.golang.org/grpc"
)

type GrpcClient struct {
	conn *grpc.ClientConn
	proto.RecordServiceClient
}

func NewGrpcClient(conn *grpc.ClientConn) *GrpcClient {
	return &GrpcClient{
		conn:                conn,
		RecordServiceClient: proto.NewRecordServiceClient(conn),
	}
}

func main() {
	var (
		addr   string
		err    error
		ctx    context.Context
		conn   *grpc.ClientConn
		client *GrpcClient
		in     *proto.RecordData
		res    *proto.RecordId
	)

	flag.StringVar(&addr, "addr", "localhost:5021", "grpc address")
	flag.Parse()

	defer func() {
		if conn != nil {
			conn.Close()
		}

		if err != nil {
			log.Fatal(err)
		}
	}()

	ctx = context.TODO()

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
	log.Println(">> grpc.Dial:", addr)

	client = NewGrpcClient(conn)
	// fmt.Println("~~~", client)

	in = proto.NewRecordData("collector_test", "test_biz")
	if res, err = client.Create(ctx, in); err != nil {
		return
	}

	log.Println(">> Create response:", res)
}
