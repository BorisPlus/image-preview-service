package config

type StorageServiceConfig struct {
	Storage StorageServerConfig
	Log     LogConfig
}

func NewStorageServiceConfig() *StorageServiceConfig {
	return &StorageServiceConfig{}
}
