package middleware_test

import (
	"bytes"
	"github.com/olzhasar/go-fileserver/middleware"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

type StubLogger struct {
	messages []string
}

func (s *StubLogger) Log(message string) {
	s.messages = append(s.messages, message)
}

type StubHandler struct {
	urls []string
}

func (s *StubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.urls = append(s.urls, r.URL.String())
}

func TestLoggingMiddleWare(t *testing.T) {
	logger := &StubLogger{}
	handler := &StubHandler{}

	loggedHandler := middleware.MakeLoggedHandler(handler, logger)

	urls := []string{
		"/upload",
		"/download",
		"/",
		"/download",
		"/upload",
	}

	for _, u := range urls {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, u, &bytes.Buffer{})

		loggedHandler.ServeHTTP(w, r)
	}

	if !reflect.DeepEqual(handler.urls, urls) {
		t.Errorf("Got recorded urls %q, want %q", handler.urls, urls)
	}

	if len(logger.messages) != len(urls) {
		t.Errorf("Got %d logged messages, want %d", len(logger.messages), len(urls))
	}
}
