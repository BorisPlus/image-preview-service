package config_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"

	"github.com/BorisPlus/previewer/core/config"
)

func TestHttpServiceConfig(t *testing.T) {
	testCases := []struct{ filename string }{
		{filename: "../../config/http_service.json"},
		{filename: "../../docker/config/http_service.json"},
		{filename: "../../docker.integration-test/config/http_service.json"},
	}
	var jsonConfig map[string]map[string]interface{}
	for _, testCase := range testCases {
		filename := testCase.filename
		log.Printf("Test case for the file %s\n", filename)
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = json.Unmarshal(data, &jsonConfig)
		if err != nil {
			t.Error(err.Error())
			return
		}
		HTTPServiceConfig := config.NewHTTPServiceConfig()
		err = config.LoadFromJsonFile(filename, HTTPServiceConfig)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if HTTPServiceConfig.HTTP.Host != jsonConfig["Http"]["Host"].(string) {
			t.Errorf(
				"HTTPServiceConfig.HTTP.Host is %s, but expected %s\n",
				HTTPServiceConfig.HTTP.Host,
				jsonConfig["HTTP"]["Host"].(string),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.HTTP.Host is %s as expected %s\n",
				HTTPServiceConfig.HTTP.Host,
				jsonConfig["Http"]["Host"].(string),
			)
		}
		if HTTPServiceConfig.HTTP.Port != uint16(jsonConfig["Http"]["Port"].(float64)) {
			t.Errorf(
				"HTTPServiceConfig.HTTP.Port is %d, but expected %d\n",
				HTTPServiceConfig.HTTP.Port,
				uint16(jsonConfig["Http"]["Port"].(float64)),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.HTTP.Port is %d as expected %d\n",
				HTTPServiceConfig.HTTP.Port,
				uint16(jsonConfig["Http"]["Port"].(float64)),
			)
		}
		if HTTPServiceConfig.HTTP.ReadTimeout != time.Duration(jsonConfig["Http"]["ReadTimeout"].(float64)) {
			t.Errorf(
				"HTTPServiceConfig.HTTP.ReadTimeout is %s, but expected %d\n",
				HTTPServiceConfig.HTTP.ReadTimeout,
				time.Duration(jsonConfig["Http"]["ReadTimeout"].(float64)),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.HTTP.ReadTimeout is %s as expected %v\n",
				HTTPServiceConfig.HTTP.ReadTimeout,
				time.Duration(jsonConfig["Http"]["ReadTimeout"].(float64)),
			)
		}
		if HTTPServiceConfig.HTTP.ReadHeaderTimeout != time.Duration(jsonConfig["Http"]["ReadHeaderTimeout"].(float64)) {
			t.Errorf(
				"HTTPServiceConfig.HTTP.ReadHeaderTimeout is %s, but expected %v\n",
				HTTPServiceConfig.HTTP.ReadHeaderTimeout,
				time.Duration(jsonConfig["Http"]["ReadHeaderTimeout"].(float64)),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.HTTP.ReadHeaderTimeout is %s as expected %v\n",
				HTTPServiceConfig.HTTP.ReadHeaderTimeout,
				time.Duration(jsonConfig["Http"]["ReadHeaderTimeout"].(float64)),
			)
		}
		if HTTPServiceConfig.HTTP.WriteTimeout != time.Duration(jsonConfig["Http"]["WriteTimeout"].(float64)) {
			t.Errorf(
				"HTTPServiceConfig.HTTP.WriteTimeout is %s, but expected %s\n",
				HTTPServiceConfig.HTTP.WriteTimeout,
				time.Duration(jsonConfig["Http"]["WriteTimeout"].(float64)),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.HTTP.WriteTimeout is %s as expected %v\n",
				HTTPServiceConfig.HTTP.WriteTimeout,
				time.Duration(jsonConfig["Http"]["WriteTimeout"].(float64)),
			)
		}
		if HTTPServiceConfig.HTTP.MaxHeaderBytes != int(jsonConfig["Http"]["MaxHeaderBytes"].(float64)) {
			t.Errorf(
				"HTTPServiceConfig.HTTP.MaxHeaderBytes is %d, but expected %d\n",
				HTTPServiceConfig.HTTP.MaxHeaderBytes,
				int(jsonConfig["Http"]["MaxHeaderBytes"].(float64)))
		} else {
			log.Printf(
				"HTTPServiceConfig.HTTP.MaxHeaderBytes is %d as expected %d\n",
				HTTPServiceConfig.HTTP.MaxHeaderBytes,
				int(jsonConfig["Http"]["MaxHeaderBytes"].(float64)),
			)
		}
		if HTTPServiceConfig.Storage.Host != jsonConfig["Storage"]["Host"].(string) {
			t.Errorf(
				"HTTPServiceConfig.Storage.Host is %s, but expected %s\n",
				HTTPServiceConfig.Storage.Host,
				jsonConfig["Storage"]["Host"].(string),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.Storage.Host is %s as expected %s\n",
				HTTPServiceConfig.Storage.Host,
				jsonConfig["Storage"]["Host"].(string),
			)
		}
		if HTTPServiceConfig.Storage.Port != uint16(jsonConfig["Storage"]["Port"].(float64)) {
			t.Errorf(
				"HTTPServiceConfig.Storage.Port is %d, but expected %d\n",
				HTTPServiceConfig.Storage.Port,
				uint16(jsonConfig["Storage"]["Port"].(float64)),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.Storage.Port is %d as expected %d\n",
				HTTPServiceConfig.Storage.Port,
				uint16(jsonConfig["Storage"]["Port"].(float64)),
			)
		}
		if string(HTTPServiceConfig.Log.Level) != jsonConfig["Log"]["Level"].(string) {
			t.Errorf(
				"HTTPServiceConfig.Log.Level is %v, but expected %s\n",
				HTTPServiceConfig.Log.Level,
				jsonConfig["Log"]["Level"].(string),
			)
		} else {
			log.Printf(
				"HTTPServiceConfig.Log.Level is %v as expected %s\n",
				HTTPServiceConfig.Log.Level,
				jsonConfig["Log"]["Level"].(string),
			)
		}
	}

}

func TestStorageServiceConfig(t *testing.T) {
	testCases := []struct{ filename string }{
		{filename: "../../config/storage_service.json"},
		{filename: "../../docker/config/storage_service.json"},
		{filename: "../../docker.integration-test/config/storage_service.json"},
	}
	var jsonConfig map[string]map[string]interface{}
	for _, testCase := range testCases {
		filename := testCase.filename
		log.Printf("Test case for the file %s\n", filename)
		data, err := os.ReadFile(filename)
		if err != nil {
			t.Error(err.Error())
			return
		}
		err = json.Unmarshal(data, &jsonConfig)
		if err != nil {
			t.Error(err.Error())
			return
		}
		StorageServiceConfig := config.NewStorageServiceConfig()
		err = config.LoadFromJsonFile(filename, StorageServiceConfig)
		if err != nil {
			t.Error(err.Error())
			return
		}
		if StorageServiceConfig.Storage.Host != jsonConfig["Storage"]["Host"].(string) {
			t.Errorf(
				"StorageServiceConfig.Storage.Host is %s, but expected %s\n",
				StorageServiceConfig.Storage.Host,
				jsonConfig["Storage"]["Host"].(string),
			)
		} else {
			log.Printf(
				"StorageServiceConfig.Storage.Host is %s as expected %s\n",
				StorageServiceConfig.Storage.Host,
				jsonConfig["Storage"]["Host"].(string),
			)
		}
		if StorageServiceConfig.Storage.Port != uint16(jsonConfig["Storage"]["Port"].(float64)) {
			t.Errorf(
				"StorageServiceConfig.Storage.Port is %d, but expected %d\n",
				StorageServiceConfig.Storage.Port,
				uint16(jsonConfig["Storage"]["Port"].(float64)),
			)
		} else {
			log.Printf(
				"StorageServiceConfig.Storage.Port is %d as expected %d\n",
				StorageServiceConfig.Storage.Port,
				uint16(jsonConfig["Storage"]["Port"].(float64)),
			)
		}
		if StorageServiceConfig.Storage.Capacity != int(jsonConfig["Storage"]["Capacity"].(float64)) {
			t.Errorf(
				"StorageServiceConfig.Storage.Capacity is %d, but expected %d\n",
				StorageServiceConfig.Storage.Capacity,
				int(jsonConfig["Storage"]["Capacity"].(float64)),
			)
		} else {
			log.Printf(
				"StorageServiceConfig.Storage.Capacity is %d as expected %d\n",
				StorageServiceConfig.Storage.Capacity,
				int(jsonConfig["Storage"]["Capacity"].(float64)),
			)
		}
		if string(StorageServiceConfig.Log.Level) != jsonConfig["Log"]["Level"].(string) {
			t.Errorf(
				"StorageServiceConfig.Log.Level is %v, but expected %s\n",
				StorageServiceConfig.Log.Level,
				jsonConfig["Log"]["Level"].(string),
			)
		} else {
			log.Printf(
				"StorageServiceConfig.Log.Level is %v as expected %s\n",
				StorageServiceConfig.Log.Level,
				jsonConfig["Log"]["Level"].(string),
			)
		}
	}
}
