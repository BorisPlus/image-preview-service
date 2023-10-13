package app

import (
	"context"
	"net/http"

	"github.com/BorisPlus/previewer/core/interfaces"
	"github.com/BorisPlus/previewer/core/models"
)

type ImagePreviewProviderApp struct {
	logger        interfaces.Logger
	storageClient interfaces.FrontendStorageClient
	maker         interfaces.ImagePreviewMaker
}

func NewImagePreviewProviderApp(
	logger interfaces.Logger,
	storageClient interfaces.FrontendStorageClient,
	maker interfaces.ImagePreviewMaker,
) *ImagePreviewProviderApp {
	return &ImagePreviewProviderApp{logger, storageClient, maker}
}

func (a ImagePreviewProviderApp) Select(ctx context.Context, transformation *models.Transformation) (*models.Result, error) {
	return a.storageClient.Select(ctx, transformation)
}

func (a ImagePreviewProviderApp) Insert(ctx context.Context, transformation *models.Transformation) (bool, error) {
	return a.storageClient.Insert(ctx, transformation)
}

func (a ImagePreviewProviderApp) AssignMakeImagePreview(
	ctx context.Context,
	transformation *models.Transformation,
	headers http.Header,
) (*models.Result, error) {
	return a.maker.MakeImagePreview(ctx, transformation, headers)
}
