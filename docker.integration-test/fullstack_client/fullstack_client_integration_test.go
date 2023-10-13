package fullstack_client_integration_test

import (
	"context"
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/BorisPlus/previewer/core/config"
	"github.com/BorisPlus/previewer/core/models"

	"github.com/BorisPlus/leveledlogger"
)

const PREVIEWER_LOG_PATH = "/var/log/previewer/" // env $PREVIEWER_LOG_PATH

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "config.json", "absolute path of json-format config file")
}

func TestIntegration(t *testing.T) {
	flag.Parse()
	cfg := NewFullstackClientConfig()
	err := config.LoadFromJsonFile(configFile, cfg)
	if err != nil {
		fmt.Println("Be shure you set: '-config <Absolute path of configuration file>'")
		fmt.Printf("Unable to decode into struct, %v\n", err)
		os.Exit(1)
		return
	}
	logger := leveledlogger.NewLogger(cfg.Log.Level, os.Stdout)
	logger.Debug("FullstackClientConfig: %+v", cfg)
	//
	logPath := os.Getenv("PREVIEWER_LOG_PATH")
	if logPath == "" {
		logPath = PREVIEWER_LOG_PATH
	}
	logger.Info("logPath %s", logPath)
	//
	fullstackClient := NewFullstackClient(
		context.Background(),
		*cfg,
		*logger,
		logPath,
	)
	fixtures := []*models.Transformation{
		models.NewTransformation(100, 100, fmt.Sprintf("%s/%s", cfg.NginxHTTPDSN(), "first.jpg")),
		models.NewTransformation(200, 100, fmt.Sprintf("%s/%s", cfg.NginxHTTPDSN(), "second.jpg")),
		models.NewTransformation(100, 300, fmt.Sprintf("%s/%s", cfg.NginxHTTPDSN(), "third.jpg")),
	}

	// CheckNotExists
	err = fullstackClient.CheckNotExists(fixtures...)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(2)
	}

	// CheckFailTargets
	err = fullstackClient.CheckFailTargets(fixtures...)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(3)
	}

	// TODO: CheckChangeStatus
	err = fullstackClient.CheckChangeStatus(fixtures...)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(4)
	}

	// CheckDefaultRoute
	err = fullstackClient.CheckDefaultRoute()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(5)
	}

	// Interaction by http
	err = fullstackClient.CheckTransformationRequest(
		models.NewTransformation(100, 100, fmt.Sprintf("%s/%s", cfg.NginxHTTPDSN(), "001.jpg")),
		models.NewTransformation(50, 50, fmt.Sprintf("%s/%s", cfg.NginxHTTPDSN(), "001.jpg")),
		models.NewTransformation(0, 0, fmt.Sprintf("%s/%s", cfg.NginxHTTPDSN(), "001.jpg")),
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(6)
	} else {
		logger.Info("OK, transformered '001.jpg' equals its native transformations ethalons")
	}

	// CheckProxiedHeaders
	err = fullstackClient.CheckProxiedHeaders(
		models.NewTransformation(100, 100, fmt.Sprintf("%s/%s", cfg.ExtHTTPDSN(), "image.jpg")),
		models.NewTransformation(50, 50, fmt.Sprintf("%s/%s", cfg.ExtHTTPDSN(), "image.jpg")),
		models.NewTransformation(0, 0, fmt.Sprintf("%s/%s", cfg.ExtHTTPDSN(), "image.jpg")),
	)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(7)
	} else {
		logger.Info("OK, external http-server get all proxied client logs")
	}
	// os.Exit(0)
}
