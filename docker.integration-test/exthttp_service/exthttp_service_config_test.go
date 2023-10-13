package exthttp_service_integration_test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/BorisPlus/previewer/core/config"
)

func TestExtHTTPServiceConfig(t *testing.T) {
	testCases := []struct{ filename string }{
		{filename: "../config/exthttp_service.json"},
	}
	var jsonConfig map[string]interface{}
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
		ExtHTTPConfig := config.NewHTTPServerConfig()
		err = config.LoadFromJsonFile(filename, ExtHTTPConfig)
		if err != nil {
			t.Error(err.Error())
			return
		}
		// ExternalHTTP  config.HTTPClientConfig
		if ExtHTTPConfig.Host != jsonConfig["Host"].(string) {
			t.Errorf(
				"ExtHTTPConfig.Host is %s, but expected %s\n",
				ExtHTTPConfig.Host,
				jsonConfig["Host"].(string),
			)
		} else {
			log.Printf(
				"ExtHTTPConfig.Host is %s as expected %s\n",
				ExtHTTPConfig.Host,
				jsonConfig["Host"].(string),
			)
		}
		if ExtHTTPConfig.Port != uint16(jsonConfig["Port"].(float64)) {
			t.Errorf(
				"ExtHTTPConfig.Port is %d, but expected %d\n",
				ExtHTTPConfig.Port,
				uint16(jsonConfig["Port"].(float64)),
			)
		} else {
			log.Printf(
				"ExtHTTPConfig.Port is %d as expected %d\n",
				ExtHTTPConfig.Port,
				uint16(jsonConfig["Port"].(float64)),
			)
		}
	}
}
