package e2e_test

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"regexp"
	"testing"
	"time"

	"github.com/mell62/journeyveil/internal/api"
)

var uuidV4Pattern = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

type channelWriter struct {
	writes chan []byte
}

func (w *channelWriter) Write(body []byte) (int, error) {
	copyOfBody := append([]byte(nil), body...)
	w.writes <- copyOfBody
	return len(body), nil
}

func TestRequestTracing(t *testing.T) {
	logWrites := make(chan []byte, 1)
	logger := slog.New(slog.NewJSONHandler(&channelWriter{writes: logWrites}, nil))
	server, err := api.NewServer("127.0.0.1:0", api.NewHandler(logger))
	if err != nil {
		t.Fatalf("create server: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- server.Run(ctx)
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case err := <-done:
			if err != nil {
				t.Errorf("stop server: %v", err)
			}
		case <-time.After(2 * time.Second):
			t.Error("server did not stop after cancellation")
		}
	})

	request, err := http.NewRequest(http.MethodGet, "http://"+server.Addr()+"/healthz?secret=hidden", nil)
	if err != nil {
		t.Fatalf("create request: %v", err)
	}
	request.Header.Set("X-Request-ID", "client-controlled")

	client := &http.Client{Timeout: time.Second}
	response, err := client.Do(request)
	if err != nil {
		t.Fatalf("send request: %v", err)
	}
	response.Body.Close()

	requestID := response.Header.Get("X-Request-ID")
	if !uuidV4Pattern.MatchString(requestID) {
		t.Fatalf("X-Request-ID = %q, want server-generated UUIDv4", requestID)
	}
	if requestID == "client-controlled" {
		t.Fatal("server trusted the client-supplied request ID")
	}

	var entry map[string]any
	select {
	case logEntry := <-logWrites:
		if bytes.Contains(logEntry, []byte("hidden")) {
			t.Error("access log contains query data")
		}
		if err := json.Unmarshal(logEntry, &entry); err != nil {
			t.Fatalf("decode access log: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("access log was not written")
	}

	assertLogField(t, entry, "msg", "HTTP request completed")
	assertLogField(t, entry, "request_id", requestID)
	assertLogField(t, entry, "method", http.MethodGet)
	assertLogField(t, entry, "path", "/healthz")
	assertLogField(t, entry, "status", float64(http.StatusOK))
	assertLogField(t, entry, "response_bytes", float64(len("{\"status\":\"ok\"}\n")))
}

func assertLogField(t *testing.T, entry map[string]any, key string, want any) {
	t.Helper()
	if got := entry[key]; got != want {
		t.Errorf("log field %q = %#v, want %#v", key, got, want)
	}
}
