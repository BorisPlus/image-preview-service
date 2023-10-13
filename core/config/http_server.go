package config

import (
	"time"
)

type HTTPServerConfig struct {
	Host              string
	Port              uint16
	ReadTimeout       time.Duration
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	MaxHeaderBytes    int
}

func NewHTTPServerConfig() *HTTPServerConfig {
	return &HTTPServerConfig{
		"",
		0,
		0,
		0,
		0,
		1048576,
	}
}
