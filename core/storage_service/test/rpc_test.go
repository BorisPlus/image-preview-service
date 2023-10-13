package rpc_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/BorisPlus/previewer/core/models"
	"github.com/BorisPlus/previewer/core/storage_service/client"
	"github.com/BorisPlus/previewer/core/storage_service/server"

	"github.com/BorisPlus/leveledlogger"
)

const (
	SAME_URL    = "http://image.ru"
	SAME_WIDTH  = 100
	SAME_HEIGHT = 200
	//
	CAPACITY = 1
)

var SAME_DATA = []byte(`01234`)

func TestInteraction(t *testing.T) {
	host := "localhost"
	var port uint16 = 5000
	dsn := fmt.Sprintf("%s:%d", host, port)

	ctx, cancel := context.WithCancel(context.Background())
	var once sync.Once
	defer once.Do(cancel)

	mainLogger := leveledlogger.NewLogger(leveledlogger.INFO, os.Stdout)
	cacheServer := storage_server.NewCacheServer(host, port, CAPACITY, mainLogger)

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		mainLogger.Info("GracefulStop")
		cacheServer.GracefulStop()
	}()
	//
	time.Sleep(5 * time.Second)
	//
	wg.Add(1)
	go func() {
		defer wg.Done()
		mainLogger.Info("server listening at %s", dsn)
		err := cacheServer.Start()
		if err != nil {
			t.Error(err.Error())
			once.Do(cancel)
		}
	}()
	//
	time.Sleep(5 * time.Second)
	//
	frontendClient := storage_client.FrontendClient{DSN: dsn}
	backendClient := storage_client.BackendClient{DSN: dsn}
	//
	transformation := models.NewTransformation(
		SAME_HEIGHT,
		SAME_WIDTH,
		SAME_URL,
	)
	exists, err := frontendClient.Insert(ctx, transformation)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	if exists == true {
		mainLogger.Error("FAIL. Data could not be exist in Storage.")
		t.Errorf("client update data because it exists in database, result is: %v, expected %v", exists, false)
	} else {
		mainLogger.Info("OK. FrontendClient INSERT data successfully")
	}
	//
	result, err := frontendClient.Select(ctx, transformation)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	if result.GetState() != int32(models.RAW) {
		mainLogger.Error("FAIL. INSERT data have no RAW-state.")
		t.Errorf("Already INSERT data has state: %v, expected %v", result.GetState(), models.RAW)
	} else {
		mainLogger.Info("OK. INSERT data have a RAW status")
	}
	if !bytes.Equal(result.GetData(), []byte{}) {
		mainLogger.Error("FAIL. INSERT data have no emply data.")
		t.Errorf("Already INSERT data has value: %v, expected %v", result.GetData(), []byte{})
	} else {
		mainLogger.Info("OK. INSERT data is emply")
	}
	if !models.Equal(result, models.EmptyResult) {
		mainLogger.Error("FAIL. INSERT data have no emply result.")
		t.Errorf("Already INSERT data has result: %v, expected %v", result.GetData(), []byte{})
	} else {
		mainLogger.Info("OK. INSERT data has empty result")
	}
	//
	result = models.NewResult(
		[]byte(`012345`),
		models.PROCESSING,
	)
	exists, err = backendClient.Update(ctx, models.NewTransformationWithResult(*transformation, *result))
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	if exists == false {
		mainLogger.Error("FAIL to update data")
		t.Errorf("client INSERT data as new, but identity original was set early, so result is: %v, expected %v", exists, true)
	} else {
		mainLogger.Info("OK. BackendClient UPDATE data successfully")
	}
	//
	resultSelectBackend, err := backendClient.Select(ctx, transformation)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	resultSelectFrontend, err := frontendClient.Select(ctx, transformation)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	if !bytes.Equal(resultSelectFrontend.GetData(), resultSelectBackend.GetData()) {
		mainLogger.Error("FAIL to equal Frontend and Backend same data SELECT")
		t.Errorf("Frontend SELECT: %v, Backend SELECT %v", resultSelectFrontend.GetData(), resultSelectBackend.GetData())
	} else {
		mainLogger.Info("OK. Frontend and Backend clients SELECT same data")
	}
	if resultSelectFrontend.GetState() != resultSelectBackend.GetState() {
		mainLogger.Error("FAIL to equal Frontend and Backend state of same data SELECT")
		t.Errorf("Frontend SELECT: %v, Backend SELECT %v", resultSelectFrontend.GetState(), resultSelectBackend.GetState())
	} else {
		mainLogger.Info("OK. Frontend and Backend clients SELECT same data state")
	}
	//
	lru_overflow_transformation := models.NewTransformation(
		SAME_HEIGHT,
		SAME_WIDTH+1,
		SAME_URL,
	)
	exists, err = frontendClient.Insert(ctx, lru_overflow_transformation)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	if exists == true {
		mainLogger.Error("FAIL to INSERT new data")
		t.Errorf("FrontendClient UPDATE data because it exists in database, result is: %v, expected %v", exists, true)
	} else {
		mainLogger.Info("OK. FrontendClient INSERT new data successfully, so it must stackoverflow storage")
	}
	//
	result, err = frontendClient.Select(ctx, transformation)
	if err != nil {
		mainLogger.Error(err.Error())
		return
	}
	if !models.Equal(result, models.NilResult) {
		mainLogger.Error("FAIL of SELECT overflowed old data")
		t.Errorf("FAIL. FrontendClient get not NIL data, witch must be overflowed early")
		t.Errorf("result is %v, expected %v",
			result,
			models.NilResult,
		)
	} else {
		mainLogger.Info("OK. client successful get no any early exist data, because it was be overflow")
	}
	//
	once.Do(cancel) // TODO: need codereview for same "single"-cancel
	wg.Wait()
}
