package transformations

import (
	"context"
	"strings"
	"sync"

	"fmt"
	"net/http"
	"strconv"

	"github.com/BorisPlus/previewer/core/http_service/servemux/handlers/common"
	"github.com/BorisPlus/previewer/core/interfaces"
	"github.com/BorisPlus/previewer/core/models"
	"github.com/BorisPlus/previewer/core/pixel"
)

type Narrow struct {
	Log           interfaces.Logger
	ImageProvider interfaces.ImagePreviewProvider
}

func (h Narrow) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// get transformation params
	heightString := r.Form.Get("height")
	height, err := strconv.ParseInt(heightString, 10, 32)
	if err != nil {
		h.Log.Error(err.Error())
		_ = common.ResponseInternalError(rw)
		return
	}
	widthString := r.Form.Get("width")
	width, err := strconv.ParseInt(widthString, 10, 32)
	if err != nil {
		h.Log.Error(err.Error())
		_ = common.ResponseInternalError(rw)
		return
	}
	urlString := r.Form.Get("url")
	// make transformation
	transformation := models.NewTransformation(
		int32(height),
		int32(width),
		urlString,
	)
	// select already exisis preview
	result, err := h.ImageProvider.Select(context.Background(), transformation)
	if err != nil {
		h.Log.Error(err.Error())
		err = common.ResponseInternalError(rw)
		if err != nil {
			h.Log.Error(err.Error())
		}
		return
	}
	// preview exists
	if !models.Equal(result, models.NilResult) {
		h.Log.Info("PREVIEW exists")
		h.Log.Info("GetState %s, len %d", result.GetState(), len(result.GetData()))
		errResponseResult := common.ResponseResult(rw, *result)
		if errResponseResult != nil {
			err = common.ResponseInternalError(rw)
			if err != nil {
				h.Log.Error(err.Error())
			}
		}
		return
	}
	once := sync.Once{}
	h.Log.Debug("PREVIEW does not exist")
	hj, ok := rw.(http.Hijacker)
	if ok {
		h.Log.Debug("Hijacked")
		conn, buf, err := hj.Hijack()
		if err != nil {
			h.Log.Error(err.Error())
			return
		}
		f := func() { conn.Close() }
		defer once.Do(f)
		h.Log.Debug("ResponseWhitePixel start")
		var builder strings.Builder
		builder.WriteString("HTTP/1.1 201 Created\n")
		builder.WriteString("Content-Type: image/jpeg\n")
		builder.WriteString(fmt.Sprintf("Content-Length: %s\n", pixel.WhitePixelLen))
		builder.WriteString("\n")
		builder.Write(pixel.WhitePixel)
		_, err = buf.WriteString(builder.String())
		if err != nil {
			h.Log.Error(err.Error())
			return
		}
		buf.Flush()
		h.Log.Info("ResponseWhitePixel done")
		// err = common.ResponseResult(rw, *result)
		// if err != nil {
		// 	h.Log.Error(err.Error())
		// 	return
		// }
		once.Do(f)
	} else {
		err = common.ResponseResult(rw, *result) // is models.NilResult
		if err != nil {
			h.Log.Error(err.Error())
			return
		}
	}
	_, err = h.ImageProvider.Insert(context.Background(), transformation)
	if err != nil {
		h.Log.Error(err.Error())
		return
	}
	h.Log.Info("AssignMakeImagePreview")
	_, err = h.ImageProvider.AssignMakeImagePreview(context.Background(), transformation, r.Header)
	if err != nil {
		h.Log.Error(err.Error())
		return
	}
	h.Log.Info("Assigned")
}
