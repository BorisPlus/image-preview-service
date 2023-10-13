package common

import (
	"github.com/BorisPlus/previewer/core/interfaces"
	"github.com/BorisPlus/previewer/core/models"
	storagerpcapi "github.com/BorisPlus/previewer/core/storage_service/rpc/api"
)

func ToGRPCTransformation(i interfaces.Transformation) *storagerpcapi.Transformation {
	obj := new(storagerpcapi.Transformation)
	obj.Height = i.GetHeight()
	obj.Width = i.GetWidth()
	obj.Url = i.GetUrl()
	return obj
}

func ToTransformation(i interfaces.Transformation) *models.Transformation {
	return models.NewTransformation(
		i.GetHeight(),
		i.GetWidth(),
		i.GetUrl(),
	)
}

func ToGRPCResult(i interfaces.Result) *storagerpcapi.Result {
	obj := new(storagerpcapi.Result)
	obj.Data = i.GetData()
	obj.State = storagerpcapi.Code(i.GetState())
	return obj
}

func ToResult(i *storagerpcapi.Result) *models.Result {
	if i == nil {
		return nil
	}
	return models.NewResult(
		i.GetData(),
		models.Code(int32(i.GetState())), // Could not redeclare `GetState()->Code` to `GetState()->int32` at GRPS object
	)
}
