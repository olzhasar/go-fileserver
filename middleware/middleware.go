package middleware

import (
	"fmt"
	"github.com/olzhasar/go-fileserver/loggers"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	handler http.Handler
	logger  loggers.Logger
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	l.handler.ServeHTTP(w, r)

	message := fmt.Sprintf("%s %s DURATION: %v", r.Method, r.URL, time.Since(start))

	l.logger.Log(message)
}

func MakeLoggedHandler(handler http.Handler, logger loggers.Logger) http.Handler {
	return &LoggingMiddleware{handler, logger}
}
