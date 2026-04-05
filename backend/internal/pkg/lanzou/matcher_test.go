package lanzou

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMatchIframe(t *testing.T) {
	html := `<html><body><iframe src="/file/abc123"></iframe></body></html>`
	result := MatchIframe(html)
	if result != "/file/abc123" {
		t.Errorf("expected /file/abc123, got %s", result)
	}

	// No iframe
	html = `<html><body><div>no iframe</div></body></html>`
	result = MatchIframe(html)
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}

	// Full URL
	html = `<html><body><iframe src="https://example.com/file/xyz"></iframe></body></html>`
	result = MatchIframe(html)
	if result != "https://example.com/file/xyz" {
		t.Errorf("expected full URL, got %s", result)
	}
}

func TestParsePwdAjax(t *testing.T) {
	html := `
		<script>
		$.ajax({
			url: '/ajax.php',
			type: 'POST',
			data: {action: 'load2', passwd: pwd, sign: 'abc123'}
		});
		var sign = 'abc123';
		</script>
	`

	ajaxData, err := ParsePwdAjax(html, "testpwd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ajaxData.URL != "/ajax.php" {
		t.Errorf("expected URL /ajax.php, got %s", ajaxData.URL)
	}
	if ajaxData.Data["passwd"] != "testpwd" {
		t.Errorf("expected passwd testpwd, got %s", ajaxData.Data["passwd"])
	}
	if ajaxData.Data["action"] != "load2" {
		t.Errorf("expected action load2, got %s", ajaxData.Data["action"])
	}
}

func TestParsePwdAjax_MissingURL(t *testing.T) {
	html := `<html><body>no ajax url here</body></html>`
	_, err := ParsePwdAjax(html, "pwd")
	if err == nil {
		t.Error("expected error for missing AJAX URL")
	}
}

func TestParseFolderAjax(t *testing.T) {
	html := `
		<script>
		$.ajax({
			url: '/filemoreajax.php',
			type: 'POST',
			data: {task: 47, folder_id: '123', pgs: pgs, vei: 'e11ad'}
		});
		var pgs = 1;
		</script>
	`

	ajaxData, err := ParseFolderAjax(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ajaxData.URL != "/filemoreajax.php" {
		t.Errorf("expected URL /filemoreajax.php, got %s", ajaxData.URL)
	}
	if ajaxData.Data["pg"] != "1" {
		t.Errorf("expected pg=1, got %s", ajaxData.Data["pg"])
	}
	if ajaxData.Data["vei"] != "e11ad" {
		t.Errorf("expected vei=e11ad, got %s", ajaxData.Data["vei"])
	}
}

func TestParseFolderAjax_MissingURL(t *testing.T) {
	html := `<html><body>no ajax</body></html>`
	_, err := ParseFolderAjax(html)
	if err == nil {
		t.Error("expected error for missing AJAX URL")
	}
}

func TestParseFileMoreAjax(t *testing.T) {
	html := `
		<script>
		$.ajax({
			url: '/filemore.php',
			data: {pg: pgs, folder_id: folder_id}
		});
		</script>
	`

	ajaxData, err := ParseFileMoreAjax(html)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if ajaxData.URL != "/filemore.php" {
		t.Errorf("expected URL /filemore.php, got %s", ajaxData.URL)
	}
	if ajaxData.Type != "POST" {
		t.Errorf("expected POST, got %s", ajaxData.Type)
	}
}

func TestParseFileMoreAjax_MissingURL(t *testing.T) {
	html := `<html><body>no ajax</body></html>`
	_, err := ParseFileMoreAjax(html)
	if err == nil {
		t.Error("expected error for missing AJAX URL")
	}
}

func TestBuildQueryString(t *testing.T) {
	data := map[string]string{
		"action": "load2",
		"passwd": "test",
	}

	result := BuildQueryString(data)
	// url.Values.Encode sorts keys alphabetically
	if result != "action=load2&passwd=test" {
		t.Errorf("expected action=load2&passwd=test, got %s", result)
	}

	// Empty
	result = BuildQueryString(map[string]string{})
	if result != "" {
		t.Errorf("expected empty string, got %s", result)
	}
}

func TestExtractShareTitle(t *testing.T) {
	// From title
	html := `<title>测试文件.zip - 蓝奏云</title>`
	result := ExtractShareTitle(html)
	if result != "测试文件.zip" {
		t.Errorf("expected 测试文件.zip, got %s", result)
	}

	// From span.b
	html = `<html><body><span class="b">文件.rar</span></body></html>`
	result = ExtractShareTitle(html)
	if result != "文件.rar" {
		t.Errorf("expected 文件.rar, got %s", result)
	}

	// Empty
	html = `<html><body></body></html>`
	result = ExtractShareTitle(html)
	if result != "" {
		t.Errorf("expected empty, got %s", result)
	}
}

func TestExtractShareFileSize(t *testing.T) {
	// From meta description
	html := `<meta name="description" content="文件大小：100MB|类型：压缩文件">`
	result := ExtractShareFileSize(html)
	if result != "100MB" {
		t.Errorf("expected 100MB, got %s", result)
	}

	// Unknown
	html = `<html><body>no size</body></html>`
	result = ExtractShareFileSize(html)
	if result != "未知" {
		t.Errorf("expected 未知, got %s", result)
	}
}

func TestExtractShareFileTime(t *testing.T) {
	// From n_file_infos
	html := `<span class="n_file_infos">2024-01-15</span>`
	result := ExtractShareFileTime(html)
	if result != "2024-01-15" {
		t.Errorf("expected 2024-01-15, got %s", result)
	}

	// From 上传时间
	html = `<span>上传时间：</span>2024-02-20<br>`
	result = ExtractShareFileTime(html)
	if result != "2024-02-20" {
		t.Errorf("expected 2024-02-20, got %s", result)
	}

	// Empty
	html = `<html><body>no time</body></html>`
	result = ExtractShareFileTime(html)
	if result != "" {
		t.Errorf("expected empty, got %s", result)
	}
}

func TestDetectSharePageType(t *testing.T) {
	tests := []struct {
		html     string
		expected SharePageType
	}{
		{`<html><body><iframe src="/file/abc"></iframe></body></html>`, ShareTypeFile},
		{`<html><body><div id="passwddiv"></div></body></html>`, ShareTypePwdFile},
		{`<html><body><div id="filemore"></div></body></html>`, ShareTypeFolder},
		{`<html><body><div id="pwdload"></div></body></html>`, ShareTypePwdFolder},
		{`<html><body><div>unknown</div></body></html>`, ShareTypeUnknown},
	}

	for _, tt := range tests {
		result := DetectSharePageType(tt.html)
		if result != tt.expected {
			t.Errorf("for %q, expected %d, got %d", tt.html, tt.expected, result)
		}
	}
}

func TestParseShareFolderResponse(t *testing.T) {
	// API返回的size是数字字符串或数字
	jsonData := []byte(`{
		"zt": 1,
		"info": "success",
		"text": [
			{"id":"123","name_all":"file1.zip","size":"1024","time":"2024-01-15","icon":"zip"},
			{"id":"456","name_all":"file2.pdf","size":"2048","time":"2024-01-16","icon":"pdf"}
		]
	}`)

	files, err := ParseShareFolderResponse(jsonData)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}

	if files[0].Name != "file1.zip" {
		t.Errorf("expected file1.zip, got %s", files[0].Name)
	}
	if files[0].Size != "1024" {
		t.Errorf("expected size 1024, got %s", files[0].Size)
	}
}

func TestParseShareFolderResponse_Error(t *testing.T) {
	// Invalid JSON
	_, err := ParseShareFolderResponse([]byte(`invalid json`))
	if err == nil {
		t.Error("expected error for invalid JSON")
	}

	// Non-success status
	jsonData := []byte(`{"zt":0,"info":"error","text":[]}`)
	_, err = ParseShareFolderResponse(jsonData)
	if err == nil {
		t.Error("expected error for non-success status")
	}
}

func TestExtractShareFileInfo(t *testing.T) {
	html := `
		<title>test.zip - 蓝奏云</title>
		<meta name="description" content="文件大小：50MB|类型：zip">
		<span class="n_file_infos">2024-03-01</span>
	`

	info, err := ExtractShareFileInfo(html, "https://lanzou.com/abc", "pwd123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if info.Name != "test.zip" {
		t.Errorf("expected test.zip, got %s", info.Name)
	}
	if info.Size != "50MB" {
		t.Errorf("expected 50MB, got %s", info.Size)
	}
	if info.Time != "2024-03-01" {
		t.Errorf("expected 2024-03-01, got %s", info.Time)
	}
	if info.URL != "https://lanzou.com/abc" {
		t.Errorf("expected URL https://lanzou.com/abc, got %s", info.URL)
	}
	if info.Pwd != "pwd123" {
		t.Errorf("expected pwd123, got %s", info.Pwd)
	}
}

// Test with mock HTTP server
func TestClient_ParseShareFile_MockSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
			<html>
				<head><title>testfile.zip - 蓝奏云</title></head>
				<body>
					<meta name="description" content="文件大小：100MB|类型：zip">
					<span class="n_file_infos">2024-01-15</span>
					<iframe src="/file/abc123"></iframe>
				</body>
			</html>
		`))
	}))
	defer server.Close()

	client := NewClient("cookie")
	client.SetBaseURL(server.URL)

	info, err := ExtractShareFileInfo(`
		<title>testfile.zip - 蓝奏云</title>
		<meta name="description" content="文件大小：100MB|类型：zip">
		<span class="n_file_infos">2024-01-15</span>
	`, server.URL+"/share/abc", "")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Name != "testfile.zip" {
		t.Errorf("expected testfile.zip, got %s", info.Name)
	}
	if info.Size != "100MB" {
		t.Errorf("expected 100MB, got %s", info.Size)
	}
}
