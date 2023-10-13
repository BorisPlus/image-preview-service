package functions_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/BorisPlus/previewer/core/app/functions"

	"github.com/BorisPlus/exthttp"
	"github.com/BorisPlus/leveledlogger"
)

var (
	real_external_url      = "https://img.youtube.com/"
	SHOW_IMAGES       bool = false
)

func TestForeignDownloadByHTTP(t *testing.T) {
	headers := http.Header{}
	headers.Add("Custom-Header", "Custom-Value")
	out, err := functions.DownloadByHTTP(real_external_url, headers)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !bytes.Equal(out[0:2], []byte{255, 216}) {
		t.Error("It is not JPEG")
	}
	if !bytes.Equal(out, EthalonForeign) {
		t.Error("It is not foreign ethalon")
	}
	if SHOW_IMAGES {
		if err := os.WriteFile("000_Original_Foreign.10x10.test.jpg", out, 0444); err != nil {
			t.Error(err.Error())
		}
	}
}

func TestForeignTransform10x20(t *testing.T) {
	data, err := functions.DownloadByHTTP(real_external_url, http.Header{})
	if err != nil {
		t.Error(err.Error())
		return
	}
	out, err := functions.TransformByNearestNeighbor(10, 20, data)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !bytes.Equal(out[0:2], []byte{255, 216}) {
		t.Error("It is not JPEG")
	}
	if !bytes.Equal(out, Ethalon10x20Foreign) {
		t.Error("It is not foreign ethalon")
	}
	if SHOW_IMAGES {
		if err := os.WriteFile("001_Transform_Foreign.10x20.test.jpg", out, 0444); err != nil {
			t.Error(err.Error())
		}
	}
}

func TestForeignTransformZero(t *testing.T) {
	data, err := functions.DownloadByHTTP(real_external_url, http.Header{})
	if err != nil {
		t.Error(err.Error())
		return
	}
	out, err := functions.TransformByNearestNeighbor(0, 0, data)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !bytes.Equal(out[0:2], []byte{255, 216}) {
		t.Error("It is not JPEG")
	}
	if !bytes.Equal(out, EthalonForeign) {
		t.Error("It is not foreign ethalon")
	}
	if SHOW_IMAGES {
		if err := os.WriteFile("002_Transform_Foreign.0x0.test.jpg", out, 0444); err != nil {
			t.Error(err.Error())
		}
	}
}

func TestInternalTransformTen(t *testing.T) {
	var once sync.Once
	var internal_http_server_host string = "localhost"
	var internal_http_server_port uint16 = 8099
	url := fmt.Sprintf(
		"http://%s:%d/image.jpg",
		internal_http_server_host,
		internal_http_server_port,
	)
	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer once.Do(func() {
		ctxCancel()
		wg.Wait()
	})
	log := leveledlogger.NewLogger(leveledlogger.INFO, os.Stdout)
	externalHTTPServer := exthttp.NewInternalTestHTTPServer(
		internal_http_server_host,
		internal_http_server_port,
		log,
		"",
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Start external http-server")
		if err := externalHTTPServer.Start(); err != nil {
			log.Error(err.Error())
		}
	}()
	time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info("Stop external http-server")
		if err := externalHTTPServer.Stop(ctx); err != nil {
			log.Error("Error Stop goroutine of external http-server: " + err.Error())
		}
	}()
	//
	data, err := functions.DownloadByHTTP(url, http.Header{})
	if err != nil {
		t.Error(err.Error())
		return
	}
	if SHOW_IMAGES {
		if err := os.WriteFile("010_Original_Green.test.jpg", data, 0444); err != nil {
			log.Error(err.Error())
		}
	}
	out, err := functions.TransformByNearestNeighbor(10, 10, data)
	if err != nil {
		t.Error(err.Error())
		return
	}
	if !bytes.Equal(out[0:2], []byte{255, 216}) {
		t.Error("It is not JPEG")
	}
	if !bytes.Equal(out, Ethalon10x10Green) {
		t.Error("It is not Ethalon10x10Green")
	}
	if SHOW_IMAGES {
		if err := os.WriteFile("011_Transform_Green.10x10.test.jpg", out, 0444); err != nil {
			log.Error(err.Error())
		}
	}
	once.Do(func() {
		ctxCancel()
		wg.Wait()
	})
}

func TestFailedInternalTransformTen(t *testing.T) {
	var once sync.Once
	var internal_http_server_host string = "localhost"
	var internal_http_server_port uint16 = 8090
	url := fmt.Sprintf(
		"http://%s:%d/image2.jpg",
		internal_http_server_host,
		internal_http_server_port,
	)
	ctx, ctxCancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	defer once.Do(func() {
		ctxCancel()
		wg.Wait()
	})
	log := leveledlogger.NewLogger(leveledlogger.INFO, os.Stdout)
	externalHTTPServer := exthttp.NewInternalTestHTTPServer(
		internal_http_server_host,
		internal_http_server_port,
		log,
		"",
	)
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Start external http-server")
		if err := externalHTTPServer.Start(); err != nil {
			log.Error(err.Error())
		}
	}()
	time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		log.Info("Stop external http-server")
		if err := externalHTTPServer.Stop(ctx); err != nil {
			log.Error("Error Stop goroutine of external http-server: " + err.Error())
		}
	}()
	//
	data, err := functions.DownloadByHTTP(url, http.Header{})
	if err != nil {
		t.Error("It must be error of not valid JPEG")
		t.Error(err.Error())
		return
	}
	_, err = functions.TransformByNearestNeighbor(10, 10, data)
	if err == nil {
		t.Error("It must be error of not valid JPEG")
		return
	} else {
		log.Info("Ok. Error was expected %q", err.Error())
	}
	once.Do(func() {
		ctxCancel()
		wg.Wait()
	})
}

func TestGenerateImages(t *testing.T) {
	imageFilename := "001.jpg"
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Error("Error of get current dir")
		return
	}
	base_path := filepath.Dir(filepath.Dir(filepath.Dir(filepath.Dir(filename))))
	image_path := path.Join(base_path, "docker.integration-test", "images")
	data, err := os.ReadFile(path.Join(image_path, imageFilename))
	if err != nil {
		t.Error(err)
		return
	}
	for _, bound := range []struct{ H, W int }{
		{H: 100, W: 100},
		{H: 50, W: 100},
		{H: 50, W: 50},
		{H: 0, W: 0},
	} {
		transformed, err := functions.TransformByNearestNeighbor(bound.H, bound.W, data)
		if err != nil {
			t.Error(err)
			return
		}
		err = os.WriteFile(
			path.Join(image_path, fmt.Sprintf("transformed.%dx%d.%s", bound.H, bound.W, imageFilename)),
			transformed,
			0444,
		)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
