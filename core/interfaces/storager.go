package interfaces

import (
	"context"

	"github.com/BorisPlus/previewer/core/models" // TODO: change to interface
)

type FrontendStorageClient interface {
	Insert(context.Context, *models.Transformation) (bool, error)
	Select(context.Context, *models.Transformation) (*models.Result, error)
}

type BackendStorageClient interface {
	Update(context.Context, *models.TransformationWithResult) (bool, error)
	Select(context.Context, *models.Transformation) (*models.Result, error)
}
