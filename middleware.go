package main

import (
	"log"
	"net/http"
	"time"
)

type LoggingMiddleware struct {
	handler http.Handler
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	l.handler.ServeHTTP(w, r)

	log.Printf("%s %s DURATION: %v", r.Method, r.URL, time.Since(start))
}

func MakeLoggedHandler(handler http.Handler) *LoggingMiddleware {
	return &LoggingMiddleware{handler}
}
