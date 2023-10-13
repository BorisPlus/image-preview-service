package http_service_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/BorisPlus/previewer/core/http_service"
	"github.com/BorisPlus/previewer/core/http_service/middleware"
	"github.com/BorisPlus/previewer/core/storage_service/client"

	"github.com/BorisPlus/leveledlogger"
)

func TestMiddleware(t *testing.T) {
	var port uint16 = 8085
	ctx, ctxCancel := context.WithCancel(context.Background())
	// Server
	httpOutput := &bytes.Buffer{}
	middlewareLog := leveledlogger.NewLogger(leveledlogger.INFO, httpOutput)
	middleware.Init(middlewareLog)
	logger := loggerInstance()
	httpServer := http_service.NewHTTPServer(
		HOST,
		port,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		logger,
		AsIsImageProvider(logger, storage_client.FrontendClient{}, storage_client.BackendClient{}),
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := httpServer.Start(); err != nil { // START
			logger.Error("http server Start goroutine: %s", err.Error())
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := httpServer.Stop(ctx); err != nil { // START
			logger.Error("http server Stop goroutine: %s", err.Error())
		}
	}()
	// Wait
	time.Sleep(5 * time.Second)
	// Client
	url := fmt.Sprintf("http://%s:%d", HOST, port)
	request, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		t.Error(err.Error())
	}
	client := &http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		t.Error(err.Error())
	}
	defer resp.Body.Close()
	expectedStatusCode := http.StatusNotFound
	if resp.StatusCode != expectedStatusCode {
		t.Errorf("StatusCode must be '%d', but get '%d'\n", expectedStatusCode, resp.StatusCode)
	} else {
		logger.Info("OK. StatusCode:'%d'\n", expectedStatusCode)
	}
	//
	ctxCancel() // STOP
	wg.Wait()
}
