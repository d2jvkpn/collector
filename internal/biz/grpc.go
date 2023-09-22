package biz

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/d2jvkpn/collector/proto"

	"github.com/d2jvkpn/gotk/cloud-logging"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

func NewGSS(logger *zap.Logger, db *mongo.Database, vp *viper.Viper, otel bool) (
	gss *GrpcServiceServer, err error) {
	interceptor := logging.NewGrpcSrvLogger(logger)

	uIntes := []grpc.UnaryServerInterceptor{interceptor.Unary()}
	if otel {
		uIntes = append(uIntes, otelgrpc.UnaryServerInterceptor( /*opts ...Option*/ ))
	}

	sIntes := []grpc.StreamServerInterceptor{interceptor.Stream()}
	if otel {
		sIntes = append(sIntes, otelgrpc.StreamServerInterceptor( /*opts ...Option*/ ))
	}

	gss = &GrpcServiceServer{logger: logger, db: db}

	gss.serverOpts = []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(uIntes...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(sIntes...)),
	}

	if vp.GetBool("tls") {
		var creds credentials.TransportCredentials
		creds, err = credentials.NewServerTLSFromFile(vp.GetString("cert"), vp.GetString("key"))
		if err != nil {
			return nil, err
		}
		gss.serverOpts = append(gss.serverOpts, grpc.Creds(creds))
	}

	return gss, nil
}

type GrpcServiceServer struct {
	logger     *zap.Logger
	db         *mongo.Database
	serverOpts []grpc.ServerOption
	server     *grpc.Server
}

func (gss *GrpcServiceServer) Serve(listener net.Listener) (err error) {
	gss.server = grpc.NewServer(gss.serverOpts...)
	proto.RegisterRecordServiceServer(gss.server, gss)

	grpc_health_v1.RegisterHealthServer(gss.server, health.NewServer())

	return gss.server.Serve(listener)
}

func (gss *GrpcServiceServer) Shutdown() {
	if gss.server != nil {
		gss.server.GracefulStop()
	}
}

func (gss *GrpcServiceServer) Create(ctx context.Context, data *proto.RecordData) (
	id *proto.RecordId, err error) {

	var (
		tCtx   context.Context
		tracer trace.Tracer
		span   trace.Span
	)

	tracer = otel.Tracer("Create")
	tCtx, span = tracer.Start(ctx, "Create")
	defer func() {
		if err != nil {
			span.AddEvent(err.Error())
		}
		span.End()
	}()

	// return nil, status.Errorf(codes.Unauthenticated, "")
	createdAt := time.Now()
	at := createdAt.UTC()
	coll := fmt.Sprintf("records_%dS%d", at.Year(), (at.Month()+2)/3)
	item := RecordFromData(data, createdAt)

	result, err := gss.db.Collection(coll).InsertOne(tCtx, item)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "")
	}

	return &proto.RecordId{Id: result.InsertedID.(primitive.ObjectID).Hex()}, nil
}
