package lanzou

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClient_Task23_SetFileAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"zt":1,"info":"success"}`))
	}))
	defer server.Close()

	client := NewClient("testcookie=123")
	client.SetBaseURL(server.URL)

	resp, err := client.Task23(123, 2, "mypassword")
	if err != nil {
		t.Fatalf("Task23 returned error: %v", err)
	}
	if resp.Zt != 1 {
		t.Errorf("expected Zt=1, got %d", resp.Zt)
	}
}

func TestClient_Task23_SetFileAccess_Public(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"zt":1,"info":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("testcookie=123")
	client.SetBaseURL(server.URL)

	resp, err := client.Task23(456, 1, "")
	if err != nil {
		t.Fatalf("Task23 returned error: %v", err)
	}
	if resp.Zt != 1 {
		t.Errorf("expected Zt=1, got %d", resp.Zt)
	}
}

func TestClient_Task16_SetFolderAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"zt":1,"info":"success"}`))
	}))
	defer server.Close()

	client := NewClient("testcookie=123")
	client.SetBaseURL(server.URL)

	resp, err := client.Task16(789, 2, "folderpwd")
	if err != nil {
		t.Fatalf("Task16 returned error: %v", err)
	}
	if resp.Zt != 1 {
		t.Errorf("expected Zt=1, got %d", resp.Zt)
	}
}

func TestClient_Task23_NetworkError(t *testing.T) {
	client := NewClient("testcookie=123")
	client.SetBaseURL("http://invalid-host-that-does-not-exist.local")

	_, err := client.Task23(1, 1, "")
	if err == nil {
		t.Error("expected network error, got nil")
	}
}

func TestClient_Task16_NetworkError(t *testing.T) {
	client := NewClient("testcookie=123")
	client.SetBaseURL("http://invalid-host-that-does-not-exist.local")

	_, err := client.Task16(1, 1, "")
	if err == nil {
		t.Error("expected network error, got nil")
	}
}
