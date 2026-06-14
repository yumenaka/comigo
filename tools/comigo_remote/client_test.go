package comigo_remote

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetServerVersionUsesVersionField(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/server-info" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"Version":"v1.2.36","ServerName":"Comigo v1.2.36"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, 5)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	got, err := client.GetServerVersion()
	if err != nil {
		t.Fatalf("GetServerVersion: %v", err)
	}
	if got != "v1.2.36" {
		t.Fatalf("version = %q, want v1.2.36", got)
	}
}

func TestGetServerVersionFallsBackToServerName(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/server-info" {
			http.NotFound(w, r)
			return
		}
		_, _ = w.Write([]byte(`{"ServerName":"Comigo v1.2.35"}`))
	}))
	defer server.Close()

	client, err := NewClient(server.URL, 5)
	if err != nil {
		t.Fatalf("NewClient: %v", err)
	}
	got, err := client.GetServerVersion()
	if err != nil {
		t.Fatalf("GetServerVersion: %v", err)
	}
	if got != "v1.2.35" {
		t.Fatalf("version = %q, want v1.2.35", got)
	}
}
