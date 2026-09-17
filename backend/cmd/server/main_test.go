package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRootServesHTML(t *testing.T) {
	mux := newMux()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
	if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
		t.Fatalf("expected html content type, got %q", rr.Header().Get("Content-Type"))
	}
}

func TestHealthEndpoint(t *testing.T) {
	mux := newMux()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}
}
