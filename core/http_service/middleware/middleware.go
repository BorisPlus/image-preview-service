package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/BorisPlus/previewer/core/interfaces"
)

type Middle struct {
	logger interfaces.Logger
}

var (
	middleware *Middle
	once       sync.Once
)

func Instance() *Middle {
	if middleware == nil {
		panic("Middleware was not init by `Init(logger interfaces.Logger)`.")
	}
	return middleware
}

func Init(logger interfaces.Logger) *Middle {
	once.Do(func() {
		middleware = &Middle{}
		middleware.logger = logger
	})
	return middleware
}

func (m Middle) Listen(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		StartAt := time.Now()
		lrw := NewLoggingResponseWriter(w)
		handler.ServeHTTP(lrw, r)
		a := struct {
			StatusCode      int
			UserAgent       string
			ClientIPAddress string
			HTTPMethod      string
			HTTPVersion     string
			URLPath         string
			StartAt         time.Time
			Latency         time.Duration
		}{
			StatusCode:      lrw.StatusCode,
			UserAgent:       r.UserAgent(),
			ClientIPAddress: r.RemoteAddr,
			HTTPMethod:      r.Method,
			HTTPVersion:     r.Proto,
			URLPath:         r.URL.Path,
			StartAt:         StartAt,
			Latency:         time.Since(StartAt),
		}
		m.logger.Info("%+v", a)
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func NewLoggingResponseWriter(writer http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{writer, 0}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.StatusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
