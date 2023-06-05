package internal

import (
	"context"
	"fmt"
	"time"

	"github.com/d2jvkpn/collector/internal/biz"
	"github.com/d2jvkpn/collector/proto"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func loadGrpc(logger *zap.Logger, db *mongo.Database) {
	interceptor := proto.NewInterceptor(logger)

	uIntes := []grpc.UnaryServerInterceptor{
		interceptor.Unary(),
		otelgrpc.UnaryServerInterceptor( /*opts ...Option*/ ),
	}

	sIntes := []grpc.StreamServerInterceptor{
		interceptor.Stream(),
		otelgrpc.StreamServerInterceptor( /*opts ...Option*/ ),
	}

	iters := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(uIntes...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(sIntes...)),
	}

	_GrpcServer = grpc.NewServer(iters...)

	gss := NewGSS(db)
	proto.RegisterRecordServiceServer(_GrpcServer, gss)
}

type GrpcServiceServer struct {
	db *mongo.Database
}

func NewGSS(db *mongo.Database) *GrpcServiceServer {
	return &GrpcServiceServer{db: db}
}

func (gss *GrpcServiceServer) Create(ctx context.Context, data *proto.RecordData) (*proto.RecordId, error) {
	// return nil, status.Errorf(codes.Unauthenticated, "")
	createdAt := time.Now()
	at := createdAt.UTC()
	coll := fmt.Sprintf("records_%dS%d", at.Year(), (at.Month()+2)/3)
	item := biz.RecordFromData(data, createdAt)

	result, err := gss.db.Collection(coll).InsertOne(context.TODO(), item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	return &proto.RecordId{Id: fmt.Sprintf("%v", result.InsertedID)}, nil
}
