package storage_server

import (
	"context"
	"fmt"
	"net"

	"github.com/BorisPlus/previewer/core/interfaces"
	"github.com/BorisPlus/previewer/core/models"
	"github.com/BorisPlus/previewer/core/storage_service/common"
	storagerpcapi "github.com/BorisPlus/previewer/core/storage_service/rpc/api"

	"github.com/BorisPlus/lru"
	"google.golang.org/grpc"
)

type CacheServer struct {
	storagerpcapi.UnimplementedStorageServer
	logger  interfaces.Logger
	server  *grpc.Server
	storage lru.Cacher
	Host    string
	Port    uint16
}

func NewCacheServer(host string, port uint16, capacity int, log interfaces.Logger) *CacheServer {
	cacheServer := &CacheServer{}
	cacheServer.Host = host
	cacheServer.Port = port
	cacheServer.logger = log
	cacheServer.storage = lru.NewCache(capacity)
	return cacheServer
}

// Insert Key of Key-Value pair
//
// Retun Status of
// * false - exists and not insert
// * true - not exists and was insert
func (s *CacheServer) Insert(_ context.Context, grpcObj *storagerpcapi.Transformation) (*storagerpcapi.Status, error) {
	transformation := common.ToTransformation(grpcObj)
	key := lru.Key(transformation.Identity())
	// TODO: one transaction need
	// begin
	_, exists := s.storage.Get(key)
	if exists {
		return &storagerpcapi.Status{Exists: exists}, nil
	}
	exists = s.storage.Set(
		key,
		models.EmptyResult,
	)
	// commit
	return &storagerpcapi.Status{Exists: exists}, nil
}

func (s *CacheServer) Update(_ context.Context, grpcObj *storagerpcapi.TransformationWithResult) (*storagerpcapi.Status, error) {
	transformation := common.ToTransformation(grpcObj.GetTransformation())
	key := lru.Key(transformation.Identity())
	value := models.NewResult(
		grpcObj.Result.Data,
		models.Code(grpcObj.Result.State),
	)
	status := &storagerpcapi.Status{Exists: s.storage.Set(key, value)}
	return status, nil
}

func (s *CacheServer) Select(_ context.Context, grpcObj *storagerpcapi.Transformation) (*storagerpcapi.Result, error) {
	transformation := common.ToTransformation(grpcObj)
	key := lru.Key(transformation.Identity())
	value, truthy := s.storage.Get(key)
	if value == nil || !truthy {
		return new(storagerpcapi.Result), nil
	}
	return common.ToGRPCResult(value.(*models.Result)), nil
}

func LoggedUnaryInterceptor(logger interfaces.Logger) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Info("MIDDLEWARE %q <-- OBJECT{%s}", info.FullMethod, req)
		return handler(ctx, req)
	}
}

func (s *CacheServer) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Host, s.Port))
	if err != nil {
		return err
	}
	gRPCServer := grpc.NewServer(
		grpc.UnaryInterceptor(LoggedUnaryInterceptor(s.logger)),
	)
	storagerpcapi.RegisterStorageServer(gRPCServer, s)
	s.server = gRPCServer
	s.logger.Info("GRPCStorageServer.Start()")
	return s.server.Serve(lis)
}

func (s *CacheServer) Stop() {
	s.logger.Info("GRPCStorageServer.Stop()")
	s.server.Stop()
}

func (s *CacheServer) GracefulStop() {
	s.logger.Info("GRPCStorageServer.GracefulStop()")
	s.server.GracefulStop()
}
