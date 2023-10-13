package common

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

var faviconContent []byte = func() []byte {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return []byte{}
	}
	data, err := os.ReadFile(path.Join(filepath.Dir(filename), "favicon.ico"))
	if err != nil {
		return []byte{}
	}
	return data
}()

type FaviconHandler struct {
}

func (h FaviconHandler) ServeHTTP(rw http.ResponseWriter, _ *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "image/x-icon")
	rw.Header().Set("Content-Length", fmt.Sprintf("%d", len(faviconContent)))
	_, err := rw.Write(faviconContent) // TODO: VSCode says to check return value
	_ = err                            // but I don't want to do any log
}
