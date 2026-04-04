package lanzou

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test_cookie")
	assert.NotNil(t, client)
	assert.Equal(t, "https://pc.woozooo.com", client.baseURL)
}

func TestClientSetCookie(t *testing.T) {
	client := NewClient("")
	assert.Empty(t, client.cookie)

	client.SetCookie("new_cookie")
	assert.Equal(t, "new_cookie", client.cookie)
}

func TestExtractName(t *testing.T) {
	html := `<title>我的文件 - 蓝奏云</title>`
	name := extractName(html)
	assert.Equal(t, "我的文件", name)
}

func TestExtractNameEmpty(t *testing.T) {
	html := `<title></title>`
	name := extractName(html)
	assert.Equal(t, "", name)
}

func TestExtractDownloadURL(t *testing.T) {
	html := `<iframe src="/file/abc123" frameborder="0"></iframe>`
	url := extractDownloadURL(html, "https://wws.lanzous.com/xxx", "")
	assert.Contains(t, url, "/file/")
}

func TestExtractDownloadURLNoIframe(t *testing.T) {
	html := `<div>No iframe here</div>`
	url := extractDownloadURL(html, "https://wws.lanzous.com/xxx", "")
	assert.Empty(t, url)
}
