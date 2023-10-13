package storage_client

import (
	"context"

	"github.com/BorisPlus/previewer/core/models"
	"github.com/BorisPlus/previewer/core/storage_service/common"
	storagerpcapi "github.com/BorisPlus/previewer/core/storage_service/rpc/api"

	"github.com/BorisPlus/leveledlogger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type FrontendClient struct {
	DSN        string
	logger     leveledlogger.Logger
	grpcClient storagerpcapi.StorageClient
	connection *grpc.ClientConn
}

func NewFrontendClient(DSN string, logger leveledlogger.Logger) *FrontendClient {
	return &FrontendClient{DSN: DSN, logger: logger}
}

var localFrontendOpts = []grpc.CallOption{}

func (c *FrontendClient) connect() error {
	connection, err := grpc.Dial(c.DSN, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed on FrontendClient: %v", err)
		return err
	}
	c.connection = connection
	c.grpcClient = storagerpcapi.NewStorageClient(c.connection)
	c.logger.Info("FrontendClient connect to %q", c.DSN)
	return nil
}

func (c *FrontendClient) close() error {
	err := c.connection.Close()
	c.logger.Info("FrontendClient close connection with %q", c.DSN)
	if err != nil {
		return err
	}
	return nil
}

func (c FrontendClient) Insert(ctx context.Context, obj *models.Transformation) (bool, error) {
	err := c.connect()
	if err != nil {
		return false, err
	}
	defer c.close()
	rpcObj := common.ToGRPCTransformation(obj)
	status, err := c.grpcClient.Insert(ctx, rpcObj, localFrontendOpts...)
	c.logger.Info("FrontendClient INSERT %q with status of existing %v", obj.Identity(), status.Exists)
	return status.Exists, err
}

func (c FrontendClient) Select(ctx context.Context, obj *models.Transformation) (*models.Result, error) {
	err := c.connect()
	if err != nil {
		// return nil, err
		return models.NilResult, err
	}
	defer c.close()
	rpcObj := common.ToGRPCTransformation(obj)
	value, err := c.grpcClient.Select(ctx, rpcObj, localFrontendOpts...)
	c.logger.Info("FrontendClient SELECT %q with state %d data len(%d)", obj.Identity(), int32(value.GetState()), len(value.GetData()))
	return common.ToResult(value), err
}
