package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/BorisPlus/previewer/core/config"
	storage_server "github.com/BorisPlus/previewer/core/storage_service/server"

	"github.com/BorisPlus/leveledlogger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "absolute path of json-format config file")
}

func main() {
	if os.Args[1] == "version" {
		fmt.Println("Storage-Service@v.1")
		return
	}
	flag.Parse()
	StorageServiceConfig := config.NewStorageServiceConfig()
	err := config.LoadFromJsonFile(configFile, StorageServiceConfig)
	if err != nil {
		fmt.Println("Be shure you set: '-config <Absolute path of configuration file>'")
		fmt.Printf("Unable to decode into struct, %v\n", err)
		return
	}
	commonLog := leveledlogger.NewLogger(StorageServiceConfig.Log.Level, os.Stdout)
	commonLog.Debug("StorageServiceConfig: %+v", StorageServiceConfig)
	StorageService := storage_server.NewCacheServer(
		StorageServiceConfig.Storage.Host,
		StorageServiceConfig.Storage.Port,
		StorageServiceConfig.Storage.Capacity,
		*commonLog,
	)

	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)
	defer stop()
	wg := sync.WaitGroup{}
	var once sync.Once
	// Stop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		StorageService.Stop()
		commonLog.Info("Storage service stop")
	}()
	// Start
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := StorageService.Start(); err != nil {
			commonLog.Error("Failed to start storage server: %s", err.Error())
			once.Do(stop)
		}
	}()
	commonLog.Info("Storage service is running...")
	<-ctx.Done()
	commonLog.Info("Complex stop down was done gracefully by signal.")
	wg.Wait()
}
