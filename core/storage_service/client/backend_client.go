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

type BackendClient struct {
	DSN        string
	logger     leveledlogger.Logger
	grpcClient storagerpcapi.StorageClient
	connection *grpc.ClientConn
}

func NewBackendClient(DSN string, logger leveledlogger.Logger) *BackendClient {
	return &BackendClient{DSN: DSN, logger: logger}
}

var localBackendOpts = []grpc.CallOption{}

func (c *BackendClient) connect() error {
	connection, err := grpc.Dial(c.DSN, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.logger.Error("failed on BackendClient: %v", err)
		return err
	}
	c.connection = connection
	c.grpcClient = storagerpcapi.NewStorageClient(c.connection)
	c.logger.Info("BackendClient connect to %q", c.DSN)
	return nil
}

func (c *BackendClient) close() error {
	err := c.connection.Close()
	c.logger.Info("BackendClient close connection with %q", c.DSN)
	if err != nil {
		return err
	}
	return nil
}

func (c BackendClient) Update(ctx context.Context, obj *models.TransformationWithResult) (bool, error) {
	err := c.connect()
	if err != nil {
		return false, err
	}
	defer c.close()
	rpcObj := &storagerpcapi.TransformationWithResult{}
	rpcObj.Transformation = common.ToGRPCTransformation(obj.GetTransformation())
	rpcObj.Result = common.ToGRPCResult(obj.GetResult())
	status, err := c.grpcClient.Update(ctx, rpcObj, localBackendOpts...)
	c.logger.Info(
		"BackendClient UPDATE %q with status of existing %v and state %d data len(%d)",
		obj.GetTransformation().Identity(),
		status.Exists,
		obj.GetResult().GetState(),
		len(obj.GetResult().GetData()),
	)
	return status.Exists, err
}

func (c BackendClient) Select(ctx context.Context, obj *models.Transformation) (*models.Result, error) {
	err := c.connect()
	if err != nil {
		// return nil, err
		return models.NilResult, err
	}
	defer c.close()
	rpcObj := common.ToGRPCTransformation(obj)
	value, err := c.grpcClient.Select(ctx, rpcObj, localBackendOpts...)
	c.logger.Info("BackendClient SELECT %q with state %d data len(%d)", obj.Identity(), int32(value.GetState()), len(value.GetData()))
	return common.ToResult(value), err
}
