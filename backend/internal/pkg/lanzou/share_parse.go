package lanzou

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
)

// ShareFileInfo 分享文件信息
type ShareFileInfo struct {
	Name string `json:"name"`
	Size string `json:"size"`
	Time string `json:"time"`
	URL  string `json:"url"`
	Pwd  string `json:"pwd,omitempty"`
}

// ShareObject 分享对象
type ShareObject struct {
	Name string         `json:"name"`
	Size string         `json:"size"`
	Type string         `json:"type"` // "file" or "folder"
	List []ShareFileInfo `json:"list"`
}

// ParseShareURL 解析蓝奏云分享链接
// 返回分享的文件列表信息
func (c *Client) ParseShareURL(shareURL, pwd string) (*ShareObject, error) {
	resp, err := c.httpClient.Get(shareURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	htmlStr := string(html)

	// 检测分享类型
	isFile := strings.Contains(htmlStr, "<iframe")
	isPwdFile := strings.Contains(htmlStr, "passwddiv")
	isFolder := strings.Contains(htmlStr, "filemore")
	isPwdFolder := strings.Contains(htmlStr, "pwdload")

	// 如果页面有密码但未提供
	if (isPwdFile || isPwdFolder) && pwd == "" {
		return &ShareObject{
			Type: "file",
			Name: extractTitle(htmlStr),
		}, nil // 返回基本信息，告知需要密码
	}

	if isFile && !isPwdFile {
		// 无密码单个文件
		return c.parseShareFile(htmlStr, shareURL, "")
	} else if isPwdFile {
		// 有密码单个文件
		return c.parseShareFile(htmlStr, shareURL, pwd)
	} else if isFolder || isPwdFolder {
		// 文件夹分享
		return c.ParseShareFolder(htmlStr, shareURL, pwd)
	}

	return nil, &ParseError{
		Message: "无法解析此分享链接",
		URL:     shareURL,
	}
}

// parseShareFile 解析单个文件分享页面
func (c *Client) parseShareFile(html, shareURL, pwd string) (*ShareObject, error) {
	name := extractTitle(html)
	size := extractFileSize(html)
	timeStr := extractFileTime(html)

	result := &ShareObject{
		Name: name,
		Size: size,
		Type: "file",
		List: []ShareFileInfo{{
			Name: name,
			Size: size,
			Time: timeStr,
			URL:  shareURL,
			Pwd:  pwd,
		}},
	}

	return result, nil
}

// parseShareFolder 解析文件夹分享页面
func (c *Client) ParseShareFolder(html, shareURL, pwd string) (*ShareObject, error) {
	// 从HTML中提取文件夹文件列表
	files := extractFolderFileList(html)

	return &ShareObject{
		Name: extractTitle(html),
		Type: "folder",
		List: files,
	}, nil
}

// GetShareDownloadURL 获取分享文件的真实下载链接
// shareURL: 分享链接
// pwd: 密码(如果有)
func (c *Client) GetShareDownloadURL(shareURL, pwd string) (string, error) {
	resp, err := c.httpClient.Get(shareURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	htmlStr := string(html)

	// 检查是否需要密码
	if strings.Contains(htmlStr, "passwddiv") && pwd != "" {
		// 有密码的分享需要提交密码获取下载链接
		return c.getShareDownloadURLWithPwd(htmlStr, shareURL, pwd)
	}

	// 无密码分享，提取iframe中的下载链接
	return extractIframeDownloadURL(htmlStr, shareURL)
}

// getShareDownloadURLWithPwd 获取有密码分享文件的下载链接
func (c *Client) getShareDownloadURLWithPwd(html, shareURL, pwd string) (string, error) {
	// 从页面中提取AJAX请求参数
	ajaxURL, ajaxData, err := extractPwdAjaxParams(html, pwd)
	if err != nil {
		return "", err
	}

	// 发送AJAX请求
	reqURL := shareURL
	if !strings.HasPrefix(ajaxURL, "http") {
		reqURL = strings.TrimSuffix(shareURL, "/") + "/" + strings.TrimPrefix(ajaxURL, "/")
	} else {
		reqURL = ajaxURL
	}

	req, err := http.NewRequest("POST", reqURL, strings.NewReader(ajaxData))
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", shareURL)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析JSON响应
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	if url, ok := result["dom"].(string); ok {
		if file, ok := result["url"].(string); ok {
			return url + "/file/" + file, nil
		}
	}

	return "", &ParseError{
		Message: "无法从响应中提取下载链接",
		URL:     shareURL,
	}
}

// ===== 提取函数 =====

func extractTitle(html string) string {
	re := regexp.MustCompile(`<title>(.+?)</title>`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSuffix(match[1], " - 蓝奏云")
	}
	// 尝试其他提取方式
	re = regexp.MustCompile(`<span class="b">([^<]+)</span>`)
	match = re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func extractFileSize(html string) string {
	re := regexp.MustCompile(`文件大小[：:]([^<|]+)`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return "未知"
}

func extractFileTime(html string) string {
	re := regexp.MustCompile(`上传时间[：:]([^<]+)`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	re = regexp.MustCompile(`n_file_infos[^>]*>([^<]+)`)
	match = re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}
	return ""
}

func extractFolderFileList(html string) []ShareFileInfo {
	var files []ShareFileInfo

	// 提取文件列表项
	re := regexp.MustCompile(`<li[^>]*>.*?<a[^>]*href="([^"]*)"[^>]*>.*?<span[^>]*>([^<]+)</span>.*?</a>.*?</li>`)
	matches := re.FindAllStringSubmatch(html, -1)

	for _, match := range matches {
		if len(match) >= 3 {
			files = append(files, ShareFileInfo{
				URL:  match[1],
				Name: match[2],
			})
		}
	}

	// 如果没有匹配到，尝试其他模式
	if len(files) == 0 {
		re = regexp.MustCompile(`n_file_infos.*?>([^<]+)<`)
		matches = re.FindAllStringSubmatch(html, -1)
		for _, match := range matches {
			if len(match) >= 2 {
				files = append(files, ShareFileInfo{
					Name: match[1],
				})
			}
		}
	}

	return files
}

func extractIframeDownloadURL(html, referer string) (string, error) {
	re := regexp.MustCompile(`iframe[^>]+src=["']([^"']+)["']`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		iframeURL := match[1]
		if !strings.HasPrefix(iframeURL, "http") {
			iframeURL = "https://wws.lanzous.com" + iframeURL
		}
		return iframeURL, nil
	}

	// 尝试直接提取下载按钮链接
	re = regexp.MustCompile(`href=["'](https?://[^"']*file[^"']*)["']`)
	match = re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1], nil
	}

	return "", &ParseError{
		Message: "无法提取下载链接",
		URL:     referer,
	}
}

func extractPwdAjaxParams(html, pwd string) (url string, data string, err error) {
	// 从页面中提取AJAX请求的URL和参数
	re := regexp.MustCompile(`url\s*:\s*["']([^"']+)["']`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		url = match[1]
	}

	// 提取密码相关参数
	data = "action=load2&passwd=" + pwd

	// 如果页面有sign参数，也加上
	re = regexp.MustCompile(`sign\s*:\s*["']([^"']*)["']`)
	match = re.FindStringSubmatch(html)
	if len(match) > 1 && match[1] != "" {
		data += "&sign=" + match[1]
	}

	if url == "" {
		return "", "", &ParseError{
			Message: "无法提取AJAX参数",
		}
	}

	return url, data, nil
}

// ParseError 解析错误
type ParseError struct {
	Message string
	URL     string
}

func (e *ParseError) Error() string {
	if e.URL != "" {
		return e.Message + ": " + e.URL
	}
	return e.Message
}

// ValidateShareURL 验证是否为有效的蓝奏云分享链接
func ValidateShareURL(url string) bool {
	patterns := []string{
		`^https?://[a-z0-9-]+\.lanzou[sxio]\.com/`,
		`^https?://[a-z0-9-]+\.lanzoux\.com/`,
		`^https?://[a-z0-9-]+\.lanzouv\.com/`,
		`^https?://[a-z0-9-]+\.lanzoui\.com/`,
		`^https?://[a-z0-9-]+\.lanzouo\.com/`,
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p)
		if re.MatchString(url) {
			return true
		}
	}
	return false
}
