package interfaces

import (
	"context"
	"net/http"

	"github.com/BorisPlus/previewer/core/models" // TODO: change to interface
)

type ImagePreviewProvider interface {
	Select(context.Context, *models.Transformation) (*models.Result, error)
	Insert(context.Context, *models.Transformation) (bool, error)
	AssignMakeImagePreview(context.Context, *models.Transformation, http.Header) (*models.Result, error)
}

type ImagePreviewMaker interface {
	Select(context.Context, *models.Transformation) (*models.Result, error)
	Update(context.Context, *models.TransformationWithResult) (bool, error)
	MakeImagePreview(context.Context, *models.Transformation, http.Header) (*models.Result, error) // Download & Transform
}
