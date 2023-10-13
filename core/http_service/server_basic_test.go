package http_service_test

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/BorisPlus/previewer/core/app"
	"github.com/BorisPlus/previewer/core/http_service"
	"github.com/BorisPlus/previewer/core/http_service/middleware"
	"github.com/BorisPlus/previewer/core/interfaces"
	"github.com/BorisPlus/previewer/core/storage_service/client"

	"github.com/BorisPlus/leveledlogger"
)

const HOST = "localhost"

func asIs(_, _ int, b []byte) ([]byte, error) {
	return b, nil
}

func loggerInstance() *leveledlogger.Logger {
	return leveledlogger.NewLogger(leveledlogger.INFO, os.Stdout)
}

func AsIsImageProvider(logger interfaces.Logger,
	front interfaces.FrontendStorageClient,
	back interfaces.BackendStorageClient) *app.ImagePreviewProviderApp {
	return app.NewImagePreviewProviderApp(
		logger,
		front,
		app.NewImagePreviewMakerApp(
			logger,
			back,
			asIs,
		),
	)
}

func TestServerBasicStopNotStarted(t *testing.T) {
	var port uint16 = 8080
	log := loggerInstance()
	middleware.Init(log)
	httpServer := http_service.NewHTTPServer(
		HOST,
		port,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		AsIsImageProvider(log, storage_client.FrontendClient{}, storage_client.BackendClient{}),
	)
	err := httpServer.Stop(context.Background())
	if err != nil {
		t.Errorf(err.Error())
	}
}

func TestServerBasicStopNormally(t *testing.T) {
	var port uint16 = 8081
	log := loggerInstance()
	middleware.Init(log)
	httpServer := http_service.NewHTTPServer(
		HOST,
		port,
		10*time.Second,
		10*time.Second,
		10*time.Second,
		1<<20,
		log,
		AsIsImageProvider(log, storage_client.FrontendClient{}, storage_client.BackendClient{}),
	)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		time.Sleep(5 * time.Second)
		err := httpServer.Stop(context.Background()) // STOP
		if err != nil {
			t.Errorf(err.Error())
		}
	}()
	if err := httpServer.Start(); err != nil { // START
		log.Error("http server goroutine: %s", err.Error())
	}
	wg.Wait()
}

func TestServerBasicStopBySignalNoWait(_ *testing.T) {
	var port uint16 = 8082
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT)
	defer ctxCancel()
	logger := loggerInstance()
	middleware.Init(logger)
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
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	process, _ := os.FindProcess(int(pid))
	err := process.Signal(syscall.SIGHUP) // STOP
	if err != nil {
		logger.Error("syscall.SIGHUP error: %s", err.Error())
	}
	wg.Wait()
}

func TestServerBasicStopBySignalWithWait(_ *testing.T) {
	var port uint16 = 8083
	ctx, ctxCancel := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT)
	defer ctxCancel()
	logger := loggerInstance()
	middleware.Init(logger)
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
	time.Sleep(5 * time.Second)
	pid, _, _ := syscall.Syscall(syscall.SYS_GETPID, 0, 0, 0)
	process, _ := os.FindProcess(int(pid))
	err := process.Signal(syscall.SIGINT) // STOP
	if err != nil {
		logger.Error("syscall.SIGINT error: %s", err.Error())
	}
	wg.Wait()
}

func TestServerBasicStopByCancel(_ *testing.T) {
	var port uint16 = 8084
	ctx, ctxCancel := context.WithCancel(context.Background())
	logger := loggerInstance()
	middleware.Init(logger)
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
	time.Sleep(5 * time.Second)
	ctxCancel() // STOP
	wg.Wait()
}
