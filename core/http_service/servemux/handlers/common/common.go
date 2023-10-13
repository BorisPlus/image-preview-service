package common

import (
	"net/http"
	"strconv"

	"github.com/BorisPlus/previewer/core/models"
	"github.com/BorisPlus/previewer/core/pixel"
)

var statuses = func() map[models.Code]int {
	return map[models.Code]int{
		models.UNSPECIFIED:      http.StatusGone,
		models.RAW:              http.StatusCreated,
		models.PROCESSING:       http.StatusAccepted,
		models.PROCESSING_ERROR: http.StatusBadRequest,
		models.READY:            http.StatusOK,
		models.INTERNAL_ERROR:   http.StatusNotFound,
	}
}

func Statuses() map[models.Code]int {
	return statuses()
}

func Status(code models.Code) int {
	return statuses()[code]
}

func commonHttpHeaders(rw http.ResponseWriter) {
	rw.Header().Set("Content-Type", "image/jpeg")
	rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	rw.Header().Set("Pragma", "no-cache")
	rw.Header().Set("Expires", "0")
}

func ResponseResult(rw http.ResponseWriter, result models.Result) error {
	rw.Header().Set("Content-Length", strconv.Itoa(len(result.GetData())))
	commonHttpHeaders(rw)
	rw.WriteHeader(Status(models.Code(result.GetState())))
	if _, err := rw.Write(result.GetData()); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return err
	}
	return nil
}

func ResponseError(rw http.ResponseWriter) error {
	commonHttpHeaders(rw)
	rw.WriteHeader(http.StatusNotFound)
	_, err := rw.Write(pixel.OrangePixel)
	return err
}

func ResponseInternalError(rw http.ResponseWriter) error {
	commonHttpHeaders(rw)
	rw.WriteHeader(http.StatusInternalServerError)
	_, err := rw.Write(pixel.RedPixel)
	return err
}
