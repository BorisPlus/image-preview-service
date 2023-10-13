package exthttp_service_integration_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"testing"

	"github.com/BorisPlus/previewer/core/config"

	"github.com/BorisPlus/exthttp"
	"github.com/BorisPlus/leveledlogger"
)

var configFile string

const PREVIEWER_LOG_PATH = "/var/log/previewer/"

func init() {
	flag.StringVar(&configFile, "config", "config.json", "absolute path of json-format config file")
}

func TestIntegration(t *testing.T) {
	if os.Args[1] == "version" {
		fmt.Println("ExtHTTP-Service@v.1")
		return
	}
	flag.Parse()
	HTTPServiceConfig := config.NewHTTPServerConfig()
	err := config.LoadFromJsonFile(configFile, HTTPServiceConfig)
	if err != nil {
		fmt.Println("Be shure you set: '-config <Absolute path of configuration file>'")
		fmt.Printf("Unable to decode into struct, %v\n", err)
		return
	}
	logger := leveledlogger.NewLogger(leveledlogger.DEBUG, os.Stdout)
	logger.Info("%+v", HTTPServiceConfig)
	logPath := os.Getenv("PREVIEWER_LOG_PATH")
	if logPath == "" {
		logPath = PREVIEWER_LOG_PATH
	}
	ExtHTTPServer := exthttp.NewInternalTestHTTPServer(
		HTTPServiceConfig.Host,
		HTTPServiceConfig.Port,
		logger,
		logPath,
	)
	ctx, stop := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGTSTP)
	defer stop()
	//
	wg := sync.WaitGroup{}
	// Stop
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		_ = ExtHTTPServer.Stop(ctx)
	}()
	// Start
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = ExtHTTPServer.Start()
	}()
	// Alive
	<-ctx.Done()
	wg.Wait()
}
