package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/constellation/v1/ping", nil)

	recorder := httptest.NewRecorder()

	Ping(recorder, req)

	resp := recorder.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(recorder.Body)
	if err != nil {
		t.Errorf("Expected a nil error, but got %s", err.Error())
	}
	want := "pong"
	if string(body) != want {
		t.Errorf("Expected to get %s, but instead got %s", want, string(body))
	}
}
