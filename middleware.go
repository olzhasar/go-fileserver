package main

import (
	"fmt"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	handler http.Handler
	logger  Logger
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	l.handler.ServeHTTP(w, r)

	message := fmt.Sprintf("%s %s DURATION: %v", r.Method, r.URL, time.Since(start))

	l.logger.Log(message)
}

func MakeLoggedHandler(handler http.Handler) http.Handler {
	return &LoggingMiddleware{handler, &StdLogger{}}
}
