package common

import (
	"net/http"

	"github.com/BorisPlus/previewer/core/interfaces"
)

type DefaultHandler struct {
	Log interfaces.Logger
}

func (h DefaultHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	err := ResponseError(rw)
	if err != nil {
		h.Log.Error(err.Error())
	}
}
