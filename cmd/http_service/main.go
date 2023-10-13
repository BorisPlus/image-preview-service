package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BorisPlus/previewer/core/app"
	"github.com/BorisPlus/previewer/core/app/functions"
	"github.com/BorisPlus/previewer/core/config"
	"github.com/BorisPlus/previewer/core/http_service"
	"github.com/BorisPlus/previewer/core/http_service/middleware"
	storage_client "github.com/BorisPlus/previewer/core/storage_service/client"

	"github.com/BorisPlus/leveledlogger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "absolute path of json-format config file")
}

func main() {
	if os.Args[1] == "version" {
		fmt.Println("HTTP-Service@v.1")
		return
	}
	flag.Parse()
	HTTPServiceConfig := config.NewHTTPServiceConfig()
	err := config.LoadFromJsonFile(configFile, HTTPServiceConfig)
	if err != nil {
		fmt.Println("Be shure you set: '-config <Absolute path of configuration file>'")
		fmt.Printf("Unable to decode into struct, %v\n", err)
		return
	}
	commonLog := leveledlogger.NewLogger(HTTPServiceConfig.Log.Level, os.Stdout)
	commonLog.Info("HTTPServiceConfig: %+v", HTTPServiceConfig)
	middleware.Init(commonLog)
	commonLog.Debug("Middleware was init")
	HTTPService := http_service.NewHTTPServer(
		HTTPServiceConfig.HTTP.Host,
		HTTPServiceConfig.HTTP.Port,
		HTTPServiceConfig.HTTP.ReadTimeout,
		HTTPServiceConfig.HTTP.ReadHeaderTimeout,
		HTTPServiceConfig.HTTP.WriteTimeout,
		HTTPServiceConfig.HTTP.MaxHeaderBytes,
		*commonLog,
		app.NewImagePreviewProviderApp(
			*commonLog,
			*storage_client.NewFrontendClient(
				fmt.Sprintf("%s:%d", HTTPServiceConfig.Storage.Host, HTTPServiceConfig.Storage.Port),
				*commonLog,
			),
			app.NewImagePreviewMakerApp(
				*commonLog,
				*storage_client.NewBackendClient(
					fmt.Sprintf("%s:%d", HTTPServiceConfig.Storage.Host, HTTPServiceConfig.Storage.Port),
					*commonLog,
				),
				functions.TransformByNearestNeighbor,
			),
		),
	)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)

	var once sync.Once
	defer once.Do(stop)
	wg := sync.WaitGroup{}
	// Stop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if err := HTTPService.Stop(ctx); err != nil {
			commonLog.Error(err.Error())
		}
	}()
	// Start
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := HTTPService.Start(); err != nil {
			commonLog.Error("Failed to start HTTP server: " + err.Error())
			once.Do(stop)
		}
	}()
	commonLog.Info("Service is running...")
	<-ctx.Done()
	commonLog.Info("Complex Shutting down was done gracefully by signal.")
	wg.Wait()
	commonLog.Info("Exit.")
}
