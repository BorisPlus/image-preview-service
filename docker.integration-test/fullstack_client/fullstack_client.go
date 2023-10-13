package fullstack_client_integration_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/BorisPlus/previewer/core/http_service/servemux/handlers/common"
	"github.com/BorisPlus/previewer/core/models"
	storageClient "github.com/BorisPlus/previewer/core/storage_service/client"

	"github.com/BorisPlus/leveledlogger"
)

type FullstackClient struct {
	ctx                   context.Context
	config                FullstackClientConfig
	httpClient            *http.Client
	frontendStorageClient *storageClient.FrontendClient
	backendStorageClient  *storageClient.BackendClient
	httpLogsDir           string
	logger                leveledlogger.Logger
}

func NewFullstackClient(
	ctx context.Context,
	config FullstackClientConfig,
	logger leveledlogger.Logger,
	httpLogsDir string,
) *FullstackClient {
	fc := &FullstackClient{
		ctx:         ctx,
		config:      config,
		logger:      logger,
		httpLogsDir: httpLogsDir,
	}
	fc.httpClient = &http.Client{}
	fc.frontendStorageClient = storageClient.NewFrontendClient(fc.config.StorageDSN(), logger)
	fc.backendStorageClient = storageClient.NewBackendClient(fc.config.StorageDSN(), logger)
	return fc
}

func (fc *FullstackClient) CheckNotExists(fixture ...*models.Transformation) error {
	fc.logger.Info("CheckNotExists()")
	for _, transformation := range fixture {
		// frontend storage client
		result, err := fc.frontendStorageClient.Select(fc.ctx, transformation)
		if err != nil {
			return err
		}
		if !models.Equal(result, models.NilResult) {
			return fmt.Errorf(
				"FRONTEND. transformation %q must has NilResult, but exists with %+v",
				transformation.Identity(),
				result,
			)
		}
		// backend storage client
		result, err = fc.backendStorageClient.Select(fc.ctx, transformation)
		if err != nil {
			return err
		}
		if !models.Equal(result, models.NilResult) {
			return fmt.Errorf(
				"BACKEND. transformation %q must has NilResult, but exists with %+v",
				transformation.Identity(),
				result,
			)
		}
	}
	return nil
}

func (fc *FullstackClient) CheckFailTargets(fixture ...*models.Transformation) error {
	fc.logger.Info("CheckFailTargets()")
	for _, transformation := range fixture {
		routeTo := strings.Join(
			[]string{
				fc.config.PreviewerHTTPAddress(),
				strconv.Itoa(int(transformation.GetHeight())),
				strconv.Itoa(int(transformation.GetWidth())),
				transformation.GetUrl(),
			}, "/")
		fc.logger.Info("Request to %q", routeTo)
		request, err := http.NewRequestWithContext(fc.ctx, http.MethodGet, routeTo, nil)
		if err != nil {
			fc.logger.Error(err.Error())
			return err
		}
		response, err := fc.httpClient.Do(request)
		if err != nil {
			fc.logger.Error(err.Error())
			return err
		}
		if response.StatusCode != http.StatusCreated {
			message := fmt.Sprintf(
				"FAIL: error http-status response get %d, expected: %d",
				response.StatusCode,
				http.StatusCreated,
			)
			fc.logger.Error(message)
			response.Body.Close()
			return fmt.Errorf(message)
		} else {
			fc.logger.Info("OK: expected http-status was response: %q", response.Status)
		}
		response.Body.Close()

		time.Sleep(5 * time.Second)
		// repeate
		response, err = fc.httpClient.Do(request)
		if err != nil {
			fc.logger.Error(err.Error())
			return err
		}
		if response.StatusCode != http.StatusBadRequest {
			message := fmt.Sprintf(
				"FAIL: error http-status response get %d, expected: %d",
				response.StatusCode,
				http.StatusBadRequest,
			)
			fc.logger.Error(message)
			response.Body.Close()
			return fmt.Errorf(message)
		} else {
			fc.logger.Info("OK: expected fail http-status was response: %q", response.Status)
		}
		response.Body.Close()
	}
	return nil
}

// TODO:
func (fc *FullstackClient) CheckChangeStatus(fixture ...*models.Transformation) error {
	fc.logger.Info("CheckChangeStatus()")
	changeTo := models.NilResult
	for stateCode := range common.Statuses() {
		if stateCode == 0 {
			continue
		}
		changeTo.SetState(stateCode)
		for _, transformation := range fixture {
			// Change on backend storage client
			indicator, err := fc.backendStorageClient.Update(
				fc.ctx,
				models.NewTransformationWithResult(
					*transformation,
					*changeTo,
				),
			)
			if err != nil {
				return err
			}
			if !indicator {
				return fmt.Errorf(
					"BACKEND. transformation %+v not found",
					transformation.Identity(),
				)
			}
			// Check by frontend storage client
			selected, err := fc.frontendStorageClient.Select(fc.ctx, transformation)
			if err != nil {
				return err
			}
			if !models.Equal(selected, changeTo) {
				return fmt.Errorf(
					"FRONTEND. transformation %q must has %d with RAW-state, but has %d",
					transformation.Identity(),
					changeTo.GetState(),
					selected.GetState(),
				)
			}
			time.Sleep(5 * time.Second)
			// Check by http client
			routeTo := strings.Join(
				[]string{
					fc.config.PreviewerHTTPAddress(),
					strconv.Itoa(int(transformation.GetHeight())),
					strconv.Itoa(int(transformation.GetWidth())),
					transformation.GetUrl(),
				}, "/")
			fc.logger.Info(routeTo)
			request, err := http.NewRequestWithContext(fc.ctx, http.MethodGet, routeTo, nil)
			if err != nil {
				return err
			}
			response, err := fc.httpClient.Do(request)
			if err != nil {
				return err
			}
			response.Body.Close()
			if response.StatusCode != common.Status(stateCode) {
				message := fmt.Sprintf(
					"FAIL: error http-status response get %d, expected: %d",
					response.StatusCode,
					http.StatusAccepted,
				)
				return fmt.Errorf(message)
			} else {
				fc.logger.Info("OK: expected fail http-status was response %q", response.Status)
			}
		}
	}
	return nil
}

func (fc *FullstackClient) CheckDefaultRoute() error {
	fc.logger.Info("CheckDefaultRoute()")
	defaultRoute := strings.Join([]string{fc.config.PreviewerHTTPAddress(), "default"}, "/")
	request, err := http.NewRequestWithContext(fc.ctx, http.MethodGet, defaultRoute, nil)
	if err != nil {
		return err
	}
	response, err := fc.httpClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusBadGateway {
		message := fmt.Sprintf(
			"FAIL: error http-status response get %q, expected: %d",
			response.Status,
			http.StatusBadGateway,
		)
		return fmt.Errorf(message)
	}
	fc.logger.Info("OK. default route have http-status %q", response.Status)
	return nil
}

type TransformationTestCase struct {
	TransformationHeigth    int
	TransformationWidth     int
	TransformationImage     string
	TransformedImageEthalon string
}

func NewTransformationTestCase(
	height int,
	width int,
	imageUrl string,
	ethalonResultImageUrl string,
) *TransformationTestCase {
	return &TransformationTestCase{
		TransformationHeigth:    height,
		TransformationWidth:     width,
		TransformationImage:     imageUrl,
		TransformedImageEthalon: ethalonResultImageUrl,
	}
}

func (fc *FullstackClient) CheckTransformationRequest(fixture ...*models.Transformation) error {
	fc.logger.Info("CheckTransformationRequest()")
	retryCount := 5
	for _, transformation := range fixture {
		// External HTTP page: http://localhost:8099/transformed.100x100.100.jpg
		fc.logger.Info("CheckTransformationRequest() of %q", transformation.Identity())
		//
		tokens := strings.Split(transformation.GetUrl(), "/")
		imageFilename := tokens[len(tokens)-1]
		transformedImageRoute := fmt.Sprintf(
			"http://%s:%d/transformed.%dx%d.%s",
			fc.config.NginxHTTP.Host,
			fc.config.NginxHTTP.Port,
			transformation.GetHeight(),
			transformation.GetWidth(),
			imageFilename,
		)
		fc.logger.Info("transformedImageRoute %q", transformedImageRoute)
		request, err := http.NewRequestWithContext(fc.ctx, http.MethodGet, transformedImageRoute, nil)
		if err != nil {
			return err
		}
		response, err := fc.httpClient.Do(request)
		if err != nil {
			return err
		}
		transformedImageEthalon, err := io.ReadAll(response.Body)
		if err != nil {
			response.Body.Close()
			return err
		}
		response.Body.Close()
		// Previewer HTTP page: http://localhost:8000/100/100/localhost:8099/100.jpg
		previewerRoute := fmt.Sprintf(
			"http://%s:%d/%d/%d/%s",
			fc.config.PreviewerHTTP.Host,
			fc.config.PreviewerHTTP.Port,
			transformation.GetHeight(),
			transformation.GetWidth(),
			transformation.GetUrl(),
		)
		fc.logger.Info("previewerRoute %q", previewerRoute)
		request, err = http.NewRequestWithContext(fc.ctx, http.MethodGet, previewerRoute, nil)
		if err != nil {
			return err
		}
		retryLimit := retryCount
		transformedImage := []byte{}
		for retryLimit > 0 {
			fc.logger.Info("retryLimit %d previewerRoute %q", retryLimit, previewerRoute)
			response, err := fc.httpClient.Do(request)
			if err != nil {
				return err
			}
			fc.logger.Info("retryLimit %d previewerRoute %q status %q", retryLimit, previewerRoute, response.Status)
			if response.StatusCode == http.StatusOK {
				transformedImage, err = io.ReadAll(response.Body)
				if err != nil {
					response.Body.Close()
					return err
				}
				response.Body.Close()
				break
			}
			response.Body.Close()
			retryLimit -= 1
			time.Sleep(5 * time.Second)
		}
		if retryLimit == 0 {
			return fmt.Errorf("FAIL: error http request %q", previewerRoute)
		}
		//
		fc.logger.Info("transformedImage len() %d previewerRoute %q", len(transformedImage), previewerRoute)
		if !bytes.Equal(transformedImage, transformedImageEthalon) {
			message := fmt.Sprintf(
				"FAIL: transformedEthalon %q len()=%d not equal native transformed %q len()=%d",
				transformedImageRoute,
				len(transformedImageRoute),
				previewerRoute,
				len(previewerRoute),
			)
			// TODO: save broken images locally to demonstration
			return fmt.Errorf(message)
		}
		// frontend storage client
		result, err := fc.frontendStorageClient.Select(fc.ctx, transformation)
		if err != nil {
			return err
		}
		if !bytes.Equal(result.GetData(), transformedImageEthalon) {
			return fmt.Errorf(
				"FRONTEND. transformation %q not equal ethalon native transformed",
				transformation.Identity(),
			)
		}
		// backend storage client
		result, err = fc.backendStorageClient.Select(fc.ctx, transformation)
		if err != nil {
			return err
		}
		if !bytes.Equal(result.GetData(), transformedImageEthalon) {
			return fmt.Errorf(
				"BACKEND. transformation %q not equal ethalon native transformed",
				transformation.Identity(),
			)
		}
		fc.logger.Info("CheckTransformationRequest() of %q OK", transformation.Identity())
	}
	return nil
}

func (fc *FullstackClient) CheckProxiedHeaders(fixture ...*models.Transformation) error {
	fc.logger.Info("CheckProxiedHeaders()")
	key := "Custom-Header-From-Client"
	for _, transformation := range fixture {
		randValue := strconv.Itoa(rand.Intn(100))
		// External HTTP page: http://localhost:8099/transformed.100x100.100.jpg
		fc.logger.Info("CheckProxiedHeaders() of %q", transformation.Identity())
		//
		transformedImageRoute := fmt.Sprintf(
			"http://%s:%d/%d/%d/%s",
			fc.config.PreviewerHTTP.Host,
			fc.config.PreviewerHTTP.Port,
			transformation.GetHeight(),
			transformation.GetWidth(),
			transformation.GetUrl(),
		)
		fc.logger.Info("transformedImageRoute %q", transformedImageRoute)
		request, err := http.NewRequestWithContext(fc.ctx, http.MethodGet, transformedImageRoute, nil)
		if err != nil {
			return err
		}
		request.Header.Add(key, randValue)
		_, err = fc.httpClient.Do(request)
		if err != nil {
			return err
		}
		// http headers data
		fc.logger.Info("sleep...")
		time.Sleep(5 * time.Second)
		// outside http log-data
		filename := path.Join(fc.httpLogsDir, "headers.json")
		fc.logger.Info("path %q", fc.httpLogsDir)
		file, err := os.Open(filename)
		if err != nil {
			return err
		}
		//
		dataLog, err := io.ReadAll(file)
		if err != nil {
			file.Close()
			return err
		}
		//
		var jsonLog map[string][]string
		err = json.Unmarshal(dataLog, &jsonLog)
		if err != nil {
			file.Close()
			return err
		}
		if !reflect.DeepEqual(jsonLog[key], []string{randValue}) {
			fc.logger.Error("FAIL: get headers %s\n", jsonLog)
			fc.logger.Error("FAIL: but expect %s\n", request.Header)
			file.Close()
			continue
		} else {
			fc.logger.Info("OK: Log contain Headers")
			fc.logger.Info("OK: Headers send to proxy %+v", request.Header)
			fc.logger.Info("OK: Headers get by target %+v", jsonLog)
		}
		fc.logger.Info("CheckProxiedHeaders() of %q OK", transformation.Identity())
		file.Close()
		err = os.Remove(filename)
		if err != nil {
			fc.logger.Error(err.Error())
		}
	}
	return nil
}
