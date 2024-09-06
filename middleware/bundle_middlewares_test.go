package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBundleMiddlewares(t *testing.T) {
	tests := []struct {
		name           string
		middlewares    []Middleware
		expectedStatus int
	}{
		{
			name:           "No middlewares",
			middlewares:    []Middleware{},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Multiple middlewares",
			middlewares: []Middleware{
				func(next http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Test-1", "true")
						next.ServeHTTP(w, r)
					}
				},
				func(next http.HandlerFunc) http.HandlerFunc {
					return func(w http.ResponseWriter, r *http.Request) {
						w.Header().Set("X-Test-2", "true")
						next.ServeHTTP(w, r)
					}
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}

			bundled := BundleMiddlewares(handler, tt.middlewares...)

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			recorder := httptest.NewRecorder()

			bundled.ServeHTTP(recorder, req)

			if recorder.Code != tt.expectedStatus {
				t.Errorf("Expected status code %d, but got %d", tt.expectedStatus, recorder.Code)
			}

			if len(tt.middlewares) > 0 {
				if recorder.Header().Get("X-Test-1") != "true" {
					t.Error("Expected X-Test-1 header to be set")
				}
				if recorder.Header().Get("X-Test-2") != "true" {
					t.Error("Expected X-Test-2 header to be set")
				}
			}
		})
	}
}

func TestBundleMiddlewaresOrder(t *testing.T) {
	order := []string{}
	middleware1 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "middleware1")
			next.ServeHTTP(w, r)
		}
	}
	middleware2 := func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			order = append(order, "middleware2")
			next.ServeHTTP(w, r)
		}
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		order = append(order, "handler")
	}

	bundled := BundleMiddlewares(handler, middleware1, middleware2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	recorder := httptest.NewRecorder()

	bundled.ServeHTTP(recorder, req)

	expectedOrder := []string{"middleware1", "middleware2", "handler"}
	for i, v := range expectedOrder {
		if order[i] != v {
			t.Errorf("Expected %s at position %d, but got %s", v, i, order[i])
		}
	}
}
