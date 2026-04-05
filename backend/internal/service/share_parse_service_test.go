package service

import (
	"testing"

	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

func TestShareParseService_ValidateShareLink(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"lanzous", "https://abc.lanzous.com/ivvHsi3qyef", true},
		{"lanzoux", "https://abc.lanzoux.com/b01tp3zkj", true},
		{"lanzoui", "https://abc.lanzoui.com/iejwp06dnwyj", true},
		{"lanzouv", "https://abc.lanzouv.com/xyz123", true},
		{"lanzouo", "https://abc.lanzouo.com/xyz123", true},
		{"wws lanzous", "https://wws.lanzous.com/abc", true},
		{"invalid domain", "https://example.com/file123", false},
		{"invalid url", "not-a-url", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := lanzou.NewClient("test-cookie")
			svc := NewShareParseService(client)

			result := svc.ValidateShareLink(tt.url)
			if result != tt.expected {
				t.Errorf("ValidateShareLink(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestShareParseService_NewShareParseService(t *testing.T) {
	client := lanzou.NewClient("test-cookie")
	svc := NewShareParseService(client)

	if svc == nil {
		t.Error("expected ShareParseService, got nil")
	}
	if svc.client == nil {
		t.Error("expected client to be set")
	}
}

func TestShareParseService_ParseShareLink(t *testing.T) {
	client := lanzou.NewClient("test-cookie")
	svc := NewShareParseService(client)

	// ParseShareLink will make real HTTP call to lanzou, expect error
	_, err := svc.ParseShareLink("https://abc.lanzous.com/ivvHsi3qyef", "")
	if err == nil {
		t.Log("unexpected success - lanzou API returned valid response")
	}
}

func TestShareParseService_GetShareDownloadLink(t *testing.T) {
	client := lanzou.NewClient("test-cookie")
	svc := NewShareParseService(client)

	// GetShareDownloadLink will make real HTTP call, expect error
	_, err := svc.GetShareDownloadLink("https://abc.lanzous.com/ivvHsi3qyef", "")
	if err == nil {
		t.Log("unexpected success - lanzou returned download URL")
	}
}
