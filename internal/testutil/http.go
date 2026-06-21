package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type APIResponse struct {
	Data json.RawMessage `json:"data"`
}

func DoRequest(t *testing.T, handler http.Handler, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()

	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("failed to marshal request body: %v", err)
		}
		reader = bytes.NewReader(payload)
	}

	req := httptest.NewRequest(method, path, reader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)
	return rec
}

func DecodeResponse(t *testing.T, rec *httptest.ResponseRecorder) APIResponse {
	t.Helper()

	var response APIResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}
	return response
}

func DecodeData[T any](t *testing.T, response APIResponse) T {
	t.Helper()

	var data T
	if err := json.Unmarshal(response.Data, &data); err != nil {
		t.Fatalf("failed to decode response data: %v", err)
	}
	return data
}
