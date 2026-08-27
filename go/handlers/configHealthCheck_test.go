package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCheckHealthTarget(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := &http.Client{Timeout: time.Second}
	result := checkHealthTarget(context.Background(), client, healthTarget{
		id: "test", name: "Test", url: server.URL,
	})
	if result.Status != "up" || result.HTTPCode != http.StatusNoContent {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestCheckHealthTargetNotConfigured(t *testing.T) {
	result := checkHealthTarget(context.Background(), http.DefaultClient, healthTarget{
		id: "test", name: "Test",
	})
	if result.Status != "not-configured" {
		t.Fatalf("status = %q, want not-configured", result.Status)
	}
}

func TestCheckHealthTargetValidatesFortiWebAPIResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, _ *http.Request) {
		response.Header().Set("Content-Type", "application/json")
		_, _ = response.Write([]byte(`{"results":{"hostname":"fortiweb-lab"}}`))
	}))
	defer server.Close()

	result := checkHealthTarget(context.Background(), server.Client(), healthTarget{
		id: "fortiweb-api", name: "FortiWeb API", url: server.URL, apiTest: true, token: "test-token",
	})
	if result.Status != "up" || result.Message != "fortiweb-lab" {
		t.Fatalf("unexpected result: %#v", result)
	}
}
