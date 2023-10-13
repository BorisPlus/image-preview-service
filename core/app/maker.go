package app

import (
	"context"
	"net/http"

	"github.com/BorisPlus/previewer/core/app/functions"
	"github.com/BorisPlus/previewer/core/interfaces"
	"github.com/BorisPlus/previewer/core/models"
	"github.com/BorisPlus/previewer/core/pixel"
)

type TransformerFunc func(int, int, []byte) ([]byte, error)

type ImagePreviewMakerApp struct {
	logger        interfaces.Logger
	storageClient interfaces.BackendStorageClient
	transform     TransformerFunc
}

func NewImagePreviewMakerApp(
	logger interfaces.Logger,
	storageClient interfaces.BackendStorageClient,
	transform TransformerFunc,
) *ImagePreviewMakerApp {
	return &ImagePreviewMakerApp{logger, storageClient, transform}
}

func (a ImagePreviewMakerApp) onError(ctx context.Context, err error, transformation *models.Transformation) (*models.Result, error) {
	if err != nil {
		a.logger.Error("onError at transformation %q", transformation.Identity())
		result := models.NewResult(pixel.BlackPixel, models.PROCESSING_ERROR)
		_, errUpdate := a.Update(ctx, models.NewTransformationWithResult(
			*transformation,
			*result,
		))
		if errUpdate != nil {
			return result, errUpdate
		}
		return result, err
	}
	return models.NilResult, nil
}

func (a ImagePreviewMakerApp) download(URL string, headers http.Header) ([]byte, error) {
	defer func() { a.logger.Info("Download end %q", URL) }()
	a.logger.Info("Download start %q", URL)
	headers.Set("User-Agent", "Image-preview service")
	headers.Del("Referer")
	return functions.DownloadByHTTP(URL, headers)
}

func (a ImagePreviewMakerApp) MakeImagePreview(ctx context.Context, transformation *models.Transformation, headers http.Header) (*models.Result, error) {
	var err error
	// Accept
	result := models.NewResult(
		pixel.GrayPixel,
		models.PROCESSING,
	)
	_, err = a.Update(ctx, models.NewTransformationWithResult(
		*transformation,
		*result,
	))
	if err != nil {
		return a.onError(ctx, err, transformation)
	}
	// Download
	a.logger.Debug("download try %s", transformation.GetUrl())
	data, err := a.download(transformation.GetUrl(), headers)
	if err != nil {
		return a.onError(ctx, err, transformation)
	}
	a.logger.Debug("downloaded data len %d", len(data))
	result.SetData(pixel.GreenPixel)
	_, err = a.Update(ctx, models.NewTransformationWithResult(
		*transformation,
		*result,
	))
	if err != nil {
		return a.onError(ctx, err, transformation)
	}
	// Transform
	data, err = a.transform(int(transformation.GetHeight()), int(transformation.GetWidth()), data)
	if err != nil {
		return a.onError(ctx, err, transformation)
	}
	a.logger.Debug("transformed data len %d", len(data))
	result.SetState(models.READY)
	result.SetData(data)
	_, err = a.Update(ctx, models.NewTransformationWithResult(
		*transformation,
		*result,
	))
	if err != nil {
		return a.onError(ctx, err, transformation)
	}
	return result, nil
}

func (a ImagePreviewMakerApp) Select(ctx context.Context, transformation *models.Transformation) (*models.Result, error) {
	return a.storageClient.Select(ctx, transformation)
}

func (a ImagePreviewMakerApp) Update(ctx context.Context, twr *models.TransformationWithResult) (bool, error) {
	return a.storageClient.Update(ctx, twr)
}
