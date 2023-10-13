package config

type HTTPServiceConfig struct {
	HTTP    HTTPServerConfig
	Storage StorageClientConfig
	Log     LogConfig
}

func NewHTTPServiceConfig() *HTTPServiceConfig {
	return &HTTPServiceConfig{}
}
