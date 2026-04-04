package lanzou

import (
	"testing"
)

func TestValidateShareURL(t *testing.T) {
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
		{"wws", "https://wws.lanzous.com/abc", true},
		{"lanzoui dkbd", "https://dkbd.lanzoui.com/dkbdv7", true},
		{"invalid domain", "https://example.com/file123", false},
		{"invalid url", "not-a-url", false},
		{"empty", "", false},
		{"lanzous with subdomain", "https://xiaodao.lanzoui.com/iejwp06dnwyj", true},
		{"lanzous wrong tld", "https://abc.lanzouru.com/xyz123", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateShareURL(tt.url)
			if result != tt.expected {
				t.Errorf("ValidateShareURL(%q) = %v, want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "standard title",
			html:     "<title>test-file.zip - 蓝奏云</title>",
			expected: "test-file.zip",
		},
		{
			name:     "title with spaces",
			html:     "<title> my document.pdf - 蓝奏云</title>",
			expected: " my document.pdf",
		},
		{
			name:     "no title tag",
			html:     "<html><body>no title</body></html>",
			expected: "",
		},
		{
			name:     "empty html",
			html:     "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractTitle(tt.html)
			if result != tt.expected {
				t.Errorf("extractTitle() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractFileSize(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "standard size",
			html:     `<meta name="description" content="文件大小：2.5MB|其他信息">`,
			expected: "2.5MB",
		},
		{
			name:     "no size info",
			html:     "<html><body>no size</body></html>",
			expected: "未知",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFileSize(tt.html)
			if result != tt.expected {
				t.Errorf("extractFileSize() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractFileTime(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{
			name:     "standard time",
			html:     `<span>上传时间：2024-01-15</span>`,
			expected: "2024-01-15",
		},
		{
			name:     "no time",
			html:     "<html><body>no time</body></html>",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractFileTime(tt.html)
			if result != tt.expected {
				t.Errorf("extractFileTime() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestExtractIframeDownloadURL(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		referer     string
		expectURL   bool
		expectError bool
	}{
		{
			name:      "iframe with relative url",
			html:      `<iframe src="/file/abc123"></iframe>`,
			referer:   "https://wws.lanzous.com/abc",
			expectURL: true,
		},
		{
			name:      "iframe with absolute url",
			html:      `<iframe src="https://wws.lanzous.com/file/abc123"></iframe>`,
			referer:   "https://wws.lanzous.com/abc",
			expectURL: true,
		},
		{
			name:        "no iframe",
			html:        "<html><body>no iframe</body></html>",
			referer:     "https://wws.lanzous.com/abc",
			expectURL:   false,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := extractIframeDownloadURL(tt.html, tt.referer)

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if tt.expectURL && result == "" {
				t.Error("expected URL, got empty")
			}
		})
	}
}

func TestParseShareFile(t *testing.T) {
	client := NewClient("test-cookie")

	html := `<title>test-document.pdf - 蓝奏云</title>
<meta name="description" content="文件大小：3.2MB|其他信息">
<span>上传时间：</span>2024-01-15<br>`

	result, err := client.parseShareFile(html, "https://wws.lanzous.com/abc", "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != "test-document.pdf" {
		t.Errorf("expected name 'test-document.pdf', got '%s'", result.Name)
	}

	if result.Size != "3.2MB" {
		t.Errorf("expected size '3.2MB', got '%s'", result.Size)
	}

	if result.Type != "file" {
		t.Errorf("expected type 'file', got '%s'", result.Type)
	}

	if len(result.List) != 1 {
		t.Errorf("expected 1 list item, got %d", len(result.List))
	}
}

func TestExtractFolderFileList(t *testing.T) {
	html := `<li><a href="/file1"><span>file1.txt</span></a></li>
<li><a href="/file2"><span>file2.pdf</span></a></li>`

	result := extractFolderFileList(html)

	if len(result) != 2 {
		t.Fatalf("expected 2 files, got %d", len(result))
	}

	if result[0].Name != "file1.txt" {
		t.Errorf("expected first name 'file1.txt', got '%s'", result[0].Name)
	}

	if result[1].Name != "file2.pdf" {
		t.Errorf("expected second name 'file2.pdf', got '%s'", result[1].Name)
	}
}

func TestParseError(t *testing.T) {
	err := &ParseError{
		Message: "test error",
		URL:     "https://example.com",
	}

	expected := "test error: https://example.com"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}

	errNoURL := &ParseError{Message: "no url error"}
	if errNoURL.Error() != "no url error" {
		t.Errorf("Error() without URL = %q, want %q", errNoURL.Error(), "no url error")
	}
}
