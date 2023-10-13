package http_service_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/BorisPlus/previewer/core/app"
	"github.com/BorisPlus/previewer/core/app/functions"
	"github.com/BorisPlus/previewer/core/http_service"
	middleware "github.com/BorisPlus/previewer/core/http_service/middleware"
	"github.com/BorisPlus/previewer/core/pixel"
	storageClient "github.com/BorisPlus/previewer/core/storage_service/client"
	storageServer "github.com/BorisPlus/previewer/core/storage_service/server"

	"github.com/BorisPlus/exthttp"
)

func TestServMuxDefaultRoute(t *testing.T) {
	expected_image_bytes := pixel.OrangePixel
	expected_http_status := http.StatusNotFound
	var port uint16 = 8090
	notValidRoutePage := fmt.Sprintf("http://%s:%d/any", HOST, port)
	ctx, ctxCancel := context.WithCancel(context.Background())
	commonLog := loggerInstance()
	middleware.Init(commonLog)
	httpServer := http_service.NewHTTPServer(
		HOST,
		port,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		commonLog,
		AsIsImageProvider(commonLog, storageClient.FrontendClient{}, storageClient.BackendClient{}),
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := httpServer.Start()
		if err != nil {
			commonLog.Error("http server Start goroutine: %s", err.Error())
		}
	}()
	time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil {
			commonLog.Error("http server Stop goroutine: %s", err.Error())
		}
	}()
	//
	defer func() {
		ctxCancel()
		wg.Wait()
	}()
	//
	time.Sleep(5 * time.Second)
	// Check by service-outside client
	client := &http.Client{}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, notValidRoutePage, bytes.NewReader([]byte("")))
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %q\n", notValidRoutePage)
		return
	}
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != expected_http_status {
		t.Errorf("FAIL: error http-status response get %q, expected: %d", response.Status, expected_http_status)
		return
	} else {
		commonLog.Info("OK: expected http-status was response: %q", response.Status)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("FAIL: could not read response body: %s\n", err)
		return
	}
	if !bytes.Equal(responseData, expected_image_bytes) {
		t.Error("FAIL: error image response")
		return
	} else {
		commonLog.Info("OK: expected image was response")
	}
}

func TestServMuxWithFullInteraction(t *testing.T) {
	var service_http_server_port uint16 = 8090
	var service_cache_server_port uint16 = 5000
	service_cache_server_dsn := fmt.Sprintf("%s:%d", HOST, service_cache_server_port)
	var external_http_server_port uint16 = 8099
	service_http_server_page := fmt.Sprintf(
		"http://%s:%d/100/100/%s:%d/image.jpg",
		HOST,
		service_http_server_port,
		HOST,
		external_http_server_port,
	)
	ctx, ctxCancel := context.WithCancel(context.Background())
	commonLog := loggerInstance()
	commonLog.Info("Address %q will be check \n", service_http_server_page)
	middleware.Init(commonLog)
	httpServer := http_service.NewHTTPServer(
		HOST,
		service_http_server_port,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		*commonLog,
		app.NewImagePreviewProviderApp(
			*commonLog,
			storageClient.NewFrontendClient(service_cache_server_dsn, *commonLog),
			*app.NewImagePreviewMakerApp(
				*commonLog,
				storageClient.NewBackendClient(service_cache_server_dsn, *commonLog),
				functions.TransformByNearestNeighbor,
			),
		),
	)
	cacheServer := storageServer.NewCacheServer(
		HOST,
		service_cache_server_port,
		5,
		commonLog,
	)
	externalHTTPServer := exthttp.NewInternalTestHTTPServer(
		HOST,
		external_http_server_port,
		commonLog,
		"",
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		commonLog.Info("Start preview-service cache-server")
		err := cacheServer.Start()
		if err != nil {
			commonLog.Error(err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		commonLog.Info("Start preview-service http-server")
		err := httpServer.Start()
		if err != nil {
			commonLog.Error(err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		commonLog.Info("Start external http-server")
		err := externalHTTPServer.Start()
		if err != nil {
			commonLog.Error(err.Error())
		}
	}()
	time.Sleep(5 * time.Second)
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil { // START
			commonLog.Error("Error Stop goroutine of preview-service http-server: " + err.Error())
		}
		commonLog.Info("Stop preview-service http-server")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		cacheServer.Stop()
		commonLog.Info("Stop preview-service cache-server")
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := externalHTTPServer.Stop(ctx); err != nil { // START
			commonLog.Error("Error Stop goroutine of external http-server: " + err.Error())
		}
		commonLog.Info("Stop external http-server")
	}()
	//
	defer func() {
		ctxCancel()
		wg.Wait()
	}()
	// Check by service-outside client
	client := &http.Client{}
	// First request
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, service_http_server_page, nil)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %q\n", service_http_server_page)
		return
	}
	commonLog.Info("First request. Client init previewing of %q", service_http_server_page)
	response, err := client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusCreated {
		t.Errorf("FAIL: response StatusCode get %d, but expected: %d\n", response.StatusCode, http.StatusCreated)
		return
	} else {
		commonLog.Info("First request. Normal http-status get: %d", response.StatusCode)
	}
	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("FAIL: could not read response body: %s\n", err)
		return
	}
	commonLog.Info("First request. Client response previewing data len is: %d", len(responseData))
	response.Body.Close()
	//
	commonLog.Info("Wait for 3 sec")
	time.Sleep(5 * time.Second)
	// Second request
	request, err = http.NewRequestWithContext(ctx, http.MethodGet, service_http_server_page, nil)
	if err != nil {
		t.Errorf("FAIL: error prepare http request: %q\n", service_http_server_page)
		return
	}
	commonLog.Info("Second request. Client get previewed of %q", service_http_server_page)
	response, err = client.Do(request)
	if err != nil {
		t.Errorf("FAIL: error decode event http request: %s\n", err)
		return
	}
	if response.StatusCode != http.StatusOK {
		t.Errorf("FAIL: response StatusCode get %d, but expected: %d\n", response.StatusCode, http.StatusOK)
		return
	} else {
		commonLog.Info("Second request. Normal http-status get: %d", response.StatusCode)
	}
	responseData, err = io.ReadAll(response.Body)
	if err != nil {
		t.Errorf("FAIL: could not read response body: %s\n", err)
		return
	}
	commonLog.Info("Second request. Client response previewed data len is: %d", len(responseData))
	response.Body.Close()
	// TODO: Ethalons chech need
}
