package fullstack_client_integration_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/BorisPlus/previewer/core/config"
)

func TestFullstackClientConfig(t *testing.T) {
	testCases := []struct{ filename string }{
		{filename: "../config/fullstack_client.json"},
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
		FullstackClientConfig := NewFullstackClientConfig()
		err = config.LoadFromJsonFile(filename, FullstackClientConfig)
		if err != nil {
			t.Error(err.Error())
			return
		}
		//
		if FullstackClientConfig.PreviewerHTTP.Host != jsonConfig["PreviewerHTTP"]["Host"].(string) {
			t.Errorf(
				"FullstackClientConfig.PreviewerHTTP.Host is %s, but expected %s\n",
				FullstackClientConfig.PreviewerHTTP.Host,
				jsonConfig["PreviewerHTTP"]["Host"].(string),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.PreviewerHTTP.Host is %s as expected %s\n",
				FullstackClientConfig.PreviewerHTTP.Host,
				jsonConfig["PreviewerHTTP"]["Host"].(string),
			)
		}
		if FullstackClientConfig.PreviewerHTTP.Port != uint16(jsonConfig["PreviewerHTTP"]["Port"].(float64)) {
			t.Errorf(
				"FullstackClientConfig.PreviewerHTTP.Port is %d, but expected %d\n",
				FullstackClientConfig.PreviewerHTTP.Port,
				uint16(jsonConfig["PreviewerHTTP"]["Port"].(float64)),
			)
		} else {
			log.Printf(
				"ExtHTTPConfig.Port is %d as expected %d\n",
				FullstackClientConfig.PreviewerHTTP.Port,
				uint16(jsonConfig["PreviewerHTTP"]["Port"].(float64)),
			)
		}
		// NginxHTTP  config.HTTPClientConfig
		if FullstackClientConfig.NginxHTTP.Host != jsonConfig["NginxHTTP"]["Host"].(string) {
			t.Errorf(
				"FullstackClientConfig.NginxHTTP.Host is %s, but expected %s\n",
				FullstackClientConfig.NginxHTTP.Host,
				jsonConfig["NginxHTTP"]["Host"].(string),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.NginxHTTP.Host is %s as expected %s\n",
				FullstackClientConfig.NginxHTTP.Host,
				jsonConfig["NginxHTTP"]["Host"].(string),
			)
		}
		if FullstackClientConfig.NginxHTTP.Port != uint16(jsonConfig["NginxHTTP"]["Port"].(float64)) {
			t.Errorf(
				"FullstackClientConfig.NginxHTTP.Port is %d, but expected %d\n",
				FullstackClientConfig.NginxHTTP.Port,
				uint16(jsonConfig["NginxHTTP"]["Port"].(float64)),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.NginxHTTP.Port is %d as expected %d\n",
				FullstackClientConfig.NginxHTTP.Port,
				uint16(jsonConfig["NginxHTTP"]["Port"].(float64)),
			)
		}
		// ExtHTTP  config.HTTPClientConfig
		if FullstackClientConfig.ExtHTTP.Host != jsonConfig["ExtHTTP"]["Host"].(string) {
			t.Errorf(
				"FullstackClientConfig.ExtHTTP.Host is %s, but expected %s\n",
				FullstackClientConfig.ExtHTTP.Host,
				jsonConfig["ExtHTTP"]["Host"].(string),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.ExtHTTP.Host is %s as expected %s\n",
				FullstackClientConfig.ExtHTTP.Host,
				jsonConfig["ExtHTTP"]["Host"].(string),
			)
		}
		if FullstackClientConfig.ExtHTTP.Port != uint16(jsonConfig["ExtHTTP"]["Port"].(float64)) {
			t.Errorf(
				"FullstackClientConfig.ExtHTTP.Port is %d, but expected %d\n",
				FullstackClientConfig.ExtHTTP.Port,
				uint16(jsonConfig["ExtHTTP"]["Port"].(float64)),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.ExtHTTP.Port is %d as expected %d\n",
				FullstackClientConfig.ExtHTTP.Port,
				uint16(jsonConfig["ExtHTTP"]["Port"].(float64)),
			)
		}
		// Storage       config.StorageClientConfig
		if FullstackClientConfig.Storage.Host != jsonConfig["Storage"]["Host"].(string) {
			t.Errorf(
				"FullstackClientConfig.Storage.Host is %s, but expected %s\n",
				FullstackClientConfig.Storage.Host,
				jsonConfig["Storage"]["Host"].(string),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.Storage.Host is %s as expected %s\n",
				FullstackClientConfig.Storage.Host,
				jsonConfig["Storage"]["Host"].(string),
			)
		}
		if FullstackClientConfig.Storage.Port != uint16(jsonConfig["Storage"]["Port"].(float64)) {
			t.Errorf(
				"FullstackClientConfig.Storage.Port is %d, but expected %d\n",
				FullstackClientConfig.Storage.Port,
				uint16(jsonConfig["Storage"]["Port"].(float64)),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.Storage.Port is %d as expected %d\n",
				FullstackClientConfig.Storage.Port,
				uint16(jsonConfig["Storage"]["Port"].(float64)),
			)
		}
		// Log
		if string(FullstackClientConfig.Log.Level) != jsonConfig["Log"]["Level"].(string) {
			t.Errorf(
				"FullstackClientConfig.Log.Level is %v, but expected %s\n",
				FullstackClientConfig.Log.Level,
				jsonConfig["Log"]["Level"].(string),
			)
		} else {
			log.Printf(
				"FullstackClientConfig.Log.Level is %v as expected %s\n",
				FullstackClientConfig.Log.Level,
				jsonConfig["Log"]["Level"].(string),
			)
		}
	}
}
