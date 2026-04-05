package lanzou

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestClientBuildRequest tests the buildRequest method
func TestClientBuildRequest(t *testing.T) {
	client := NewClient("test-cookie")

	body := url.Values{"task": {"5"}, "folder_id": {"-1"}}
	req, err := client.buildRequest("POST", "doupload.php", body)

	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "POST", req.Method)
	assert.Equal(t, "https://pc.woozooo.com/doupload.php", req.URL.String())
	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36", req.Header.Get("User-Agent"))
	assert.Equal(t, "test-cookie", req.Header.Get("Cookie"))
	assert.Equal(t, "application/x-www-form-urlencoded", req.Header.Get("Content-Type"))
	assert.Equal(t, "https://pc.woozooo.com/", req.Header.Get("Referer"))
}

func TestClientBuildRequest_EmptyCookie(t *testing.T) {
	client := NewClient("")
	body := url.Values{"task": {"1"}}
	req, err := client.buildRequest("GET", "test.php", body)

	assert.NoError(t, err)
	assert.Empty(t, req.Header.Get("Cookie"))
}

func TestClientBuildRequest_InvalidMethod(t *testing.T) {
	client := NewClient("cookie")
	body := url.Values{}
	_, err := client.buildRequest("INVALID METHOD WITH SPACE\n", "test.php", body)
	assert.Error(t, err)
}

func TestClientDoRequest_Error(t *testing.T) {
	client := NewClient("cookie")
	client.baseURL = "http://localhost:1" // unreachable
	body := url.Values{"task": {"1"}}

	_, err := client.doRequest("POST", "test.php", body)
	assert.Error(t, err)
}

func TestClientPostForm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/doupload.php", r.URL.Path)
		w.Write([]byte(`{"zt":1,"info":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL
	body := url.Values{"task": {"5"}}

	resp, err := client.postForm("doupload.php", body)
	assert.NoError(t, err)
	assert.Contains(t, string(resp), `"zt":1`)
}

// Test Task5 with mock server
func TestClientTask5(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "5", r.FormValue("task"))
		assert.Equal(t, "-1", r.FormValue("folder_id"))
		w.Write([]byte(`{"zt":1,"text":[],"info":"success"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task5(-1, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
}

// Test Task47 with mock server
func TestClientTask47(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "47", r.FormValue("task"))
		w.Write([]byte(`{"zt":1,"text":[{"id":1,"name":"folder1"}],"info":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task47(-1)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
	assert.Len(t, resp.Text, 1)
}

// Test Task2 with mock server
func TestClientTask2(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "2", r.FormValue("task"))
		assert.Equal(t, "new_folder", r.FormValue("folder_name"))
		w.Write([]byte(`{"zt":1,"text":{"folder_id":123,"folder_name":"new_folder"},"info":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task2(-1, "new folder")
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
	assert.Equal(t, uint(123), resp.Text.ID)
}

// Test Task2 name sanitization
func TestClientTask2_NameSanitization(t *testing.T) {
	var receivedName string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		receivedName = r.FormValue("folder_name")
		w.Write([]byte(`{"zt":1,"text":{},"info":"ok"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	client.Task2(-1, "my folder (1)")
	assert.Equal(t, "my_folder__1_", receivedName)
}

// Test Task6 (delete file)
func TestClientTask6(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "6", r.FormValue("task"))
		assert.Equal(t, "123", r.FormValue("file_id"))
		w.Write([]byte(`{"zt":1,"info":"deleted"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task6(123)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
}

// Test Task46 (delete folder)
func TestClientTask46(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "46", r.FormValue("task"))
		w.Write([]byte(`{"zt":1,"info":"folder deleted"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task46(456)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
}

// Test Task14 (rename file)
func TestClientTask14(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "14", r.FormValue("task"))
		assert.Equal(t, "new_name.txt", r.FormValue("name"))
		w.Write([]byte(`{"zt":1,"info":"renamed"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task14(789, "new_name.txt")
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
}

// Test Task15 (move file)
func TestClientTask15(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "15", r.FormValue("task"))
		assert.Equal(t, "100", r.FormValue("file_id"))
		assert.Equal(t, "200", r.FormValue("folder_id"))
		w.Write([]byte(`{"zt":1,"info":"moved"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task15(100, 200)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
}

// Test Task48 (move folder)
func TestClientTask48(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "48", r.FormValue("task"))
		assert.Equal(t, "10", r.FormValue("folder_id"))
		assert.Equal(t, "20", r.FormValue("target_id"))
		w.Write([]byte(`{"zt":1,"info":"moved"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task48(10, 20)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
}

// Test Task39 (create share)
func TestClientTask39(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "39", r.FormValue("task"))
		assert.Equal(t, "123", r.FormValue("file_id"))
		w.Write([]byte(`{"zt":1,"text":{"url":"https://abc.lanzoui.com/x","pwd":"abc1"},"info":"shared"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task39(123, 10080)
	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
	assert.Equal(t, "https://abc.lanzoui.com/x", resp.Text.URL)
	assert.Equal(t, "abc1", resp.Text.Pwd)
}

// Test Ping
func TestClientPing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "1", r.FormValue("task"))
		w.Write([]byte(`{"zt":1}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	err := client.Ping()
	assert.NoError(t, err)
}

// Test Task5 JSON parse error
func TestClientTask5_JSONError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`not valid json`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task5(-1, 1)
	assert.NoError(t, err)
	assert.Equal(t, -1, resp.Zt) // Should fallback to -1 on parse error
}

// Test extractName
func TestExtractName_Variations(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		expected string
	}{
		{"normal", `<title>document.pdf - 蓝奏云</title>`, "document.pdf"},
		{"no suffix", `<title>file.txt</title>`, "file.txt"},
		{"empty", `<title></title>`, ""},
		{"no title", `<html><body>no title</body></html>`, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, extractName(tt.html))
		})
	}
}

// Test extractDownloadURL
func TestExtractDownloadURL_Variations(t *testing.T) {
	tests := []struct {
		name    string
		html    string
		hasURL  bool
	}{
		{"with iframe", `<iframe src="/file/abc123"></iframe>`, true},
		{"full url iframe", `<iframe src="https://example.com/f"></iframe>`, true},
		{"no iframe", `<html><body>no iframe</body></html>`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDownloadURL(tt.html, "https://referer.com", "")
			if tt.hasURL {
				assert.NotEmpty(t, result)
			} else {
				assert.Empty(t, result)
			}
		})
	}
}

// Test Task22 and Task18 (return map)
func TestClientTask22(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "22", r.FormValue("task"))
		w.Write([]byte(`{"zt":1,"file_name":"test.txt","size":1024}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task22(1)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

func TestClientTask18(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		assert.Equal(t, "18", r.FormValue("task"))
		w.Write([]byte(`{"zt":1,"folder":"test"}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	resp, err := client.Task18(1)
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}

// Test GetDownloadURL - network error path
func TestClientGetDownloadURL_NetworkError(t *testing.T) {
	client := NewClient("cookie")
	client.baseURL = "http://localhost:1" // unreachable

	_, _, err := client.GetDownloadURL("http://localhost:1/test", "")
	assert.Error(t, err)
}

// Test GetDownloadURL - mock success
func TestClientGetDownloadURL_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<title>testfile.txt - 蓝奏云</title><iframe src="/file/abc"></iframe>`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	name, downURL, err := client.GetDownloadURL(server.URL, "")
	assert.NoError(t, err)
	assert.Equal(t, "testfile.txt", name)
	assert.Contains(t, downURL, "/file/abc")
}

// Test GetDownloadURL - password required
func TestClientGetDownloadURL_PasswordRequired(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<title>file.txt</title><div id="passwddiv"></div>`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.baseURL = server.URL

	_, _, err := client.GetDownloadURL(server.URL, "")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "password")
}
