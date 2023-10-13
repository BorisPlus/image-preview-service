package fullstack_client_integration_test

import (
	"net"
	"strconv"
	"strings"

	"github.com/BorisPlus/previewer/core/config"
)

type FullstackClientConfig struct {
	PreviewerHTTP config.HTTPClientConfig
	NginxHTTP     config.HTTPClientConfig
	ExtHTTP       config.HTTPClientConfig
	Storage       config.StorageClientConfig
	Log           config.LogConfig
}

func NewFullstackClientConfig() *FullstackClientConfig {
	return &FullstackClientConfig{}
}

func (c FullstackClientConfig) PreviewerHTTPAddress() string {
	return strings.Join(
		[]string{
			"http://",
			c.PreviewerHTTP.Host,
			":",
			strconv.Itoa(int(c.PreviewerHTTP.Port)),
		}, "")
}

func (c FullstackClientConfig) NginxHTTPAddress() string {
	return strings.Join(
		[]string{
			"http://",
			c.NginxHTTP.Host,
			":",
			strconv.Itoa(int(c.NginxHTTP.Port)),
		}, "")
}

func (c FullstackClientConfig) NginxHTTPDSN() string {
	return net.JoinHostPort(c.NginxHTTP.Host, strconv.Itoa(int(c.NginxHTTP.Port)))
}

func (c FullstackClientConfig) ExtHTTPAddress() string {
	return strings.Join(
		[]string{
			"http://",
			c.ExtHTTP.Host,
			":",
			strconv.Itoa(int(c.ExtHTTP.Port)),
		}, "")
}

func (c FullstackClientConfig) ExtHTTPDSN() string {
	return net.JoinHostPort(c.ExtHTTP.Host, strconv.Itoa(int(c.ExtHTTP.Port)))
}

func (c FullstackClientConfig) StorageDSN() string {
	return net.JoinHostPort(c.Storage.Host, strconv.Itoa(int(c.Storage.Port)))
}
