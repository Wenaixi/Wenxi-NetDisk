package lanzou

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHtml5Upload_Success 测试上传成功
func TestHtml5Upload_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/html5up.php", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		w.Write([]byte(`{"zt":1,"info":"上传成功","text":[{"f_id":"123","is_newd":"https://wwn.lanzouf.com","downs":"0","icon":"txt","id":"456","name":"test.txt","size":"1024","time":"2024-01-01","onof":"0"}]}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.SetUploadBaseURL(server.URL)

	data := strings.NewReader("test file content")
	resp, err := client.Html5Upload(data, "test.txt", 17, -1)

	assert.NoError(t, err)
	assert.Equal(t, 1, resp.Zt)
	assert.Equal(t, "上传成功", resp.Info)
	assert.Len(t, resp.Text, 1)
	assert.Equal(t, "123", resp.Text[0].FID)
	assert.Equal(t, "https://wwn.lanzouf.com", resp.Text[0].IsNew)
	assert.Equal(t, "test.txt", resp.Text[0].Name)
}

// TestHtml5Upload_Failure 测试上传失败
func TestHtml5Upload_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"zt":0,"info":"文件格式不支持","text":null}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.SetUploadBaseURL(server.URL)

	data := strings.NewReader("test file content")
	resp, err := client.Html5Upload(data, "test.exe", 17, -1)

	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Zt)
	assert.Equal(t, "文件格式不支持", resp.Info)
	assert.Nil(t, resp.Text)
}

// TestHtml5Upload_NetworkError 测试网络错误
func TestHtml5Upload_NetworkError(t *testing.T) {
	client := NewClient("cookie")
	client.SetUploadBaseURL("http://127.0.0.1:1")

	data := strings.NewReader("test")
	_, err := client.Html5Upload(data, "test.txt", 4, -1)

	assert.Error(t, err)
}

// TestHtml5Upload_Headers 测试请求头
func TestHtml5Upload_Headers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Contains(t, r.Header.Get("User-Agent"), "Mozilla")
		assert.Contains(t, r.Header.Get("Referer"), "pc.woozooo.com")
		assert.NotEmpty(t, r.Header.Get("Cookie"))
		assert.Contains(t, r.Header.Get("Content-Type"), "multipart/form-data")
		w.Write([]byte(`{"zt":1,"info":"ok","text":[]}`))
	}))
	defer server.Close()

	client := NewClient("test-cookie")
	client.SetUploadBaseURL(server.URL)

	data := strings.NewReader("test")
	_, err := client.Html5Upload(data, "test.txt", 4, 123)

	assert.NoError(t, err)
}

// TestHtml5Upload_FormFields 测试表单字段
func TestHtml5Upload_FormFields(t *testing.T) {
	var capturedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 读取multipart body
		r.ParseMultipartForm(32 << 20)
		capturedBody = r.FormValue("task") + "|" + r.FormValue("vie") + "|" + r.FormValue("ve") + "|" + r.FormValue("folder_id_bb_n") + "|" + r.FormValue("name") + "|" + r.FormValue("size")
		w.Write([]byte(`{"zt":1,"info":"ok","text":[]}`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.SetUploadBaseURL(server.URL)

	data := strings.NewReader("hello world")
	_, err := client.Html5Upload(data, "hello.txt", 11, 456)

	assert.NoError(t, err)
	assert.Equal(t, "1|2|2|456|hello.txt|11", capturedBody)
}
