package http_service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/BorisPlus/previewer/core/http_service/servemux"
	"github.com/BorisPlus/previewer/core/interfaces"
)

type HTTPServer struct {
	server *http.Server
	logger interfaces.Logger
	app    interfaces.ImagePreviewProvider
}

func NewHTTPServer(
	host string,
	port uint16,
	readTimeout time.Duration, // TODO: set time.Duration default "10s" in ServerConfig `10 * time.Second`
	readHeaderTimeout time.Duration, // TODO: set default "10s" in ServerConfig `10 * time.Second`
	writeTimeout time.Duration, // TODO: set default "10s" in ServerConfig `10 * time.Second`
	maxHeaderBytes int, // TODO: set default in ServerConfig `1 << 20`
	logger interfaces.Logger,
	app interfaces.ImagePreviewProvider,
) *HTTPServer {
	mux := http.NewServeMux()
	mux.Handle("/", servemux.Handlers(logger, app))
	server := http.Server{
		Addr:              net.JoinHostPort(host, fmt.Sprint(port)),
		Handler:           mux,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
		MaxHeaderBytes:    maxHeaderBytes,
	}
	httpServer := &HTTPServer{}
	httpServer.server = &server
	httpServer.logger = logger
	httpServer.app = app
	return httpServer
}

func (s *HTTPServer) Start() error {
	s.logger.Info("HTTPServer.Start()")
	if err := s.server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		s.logger.Info("Start error: %v\n", err)
		return err
	}
	return nil
}

func (s *HTTPServer) Stop(ctx context.Context) error {
	s.logger.Info("HTTPServer.Stop()")
	err := s.server.Shutdown(ctx)
	if err != nil {
		return err
	}
	return nil
}
