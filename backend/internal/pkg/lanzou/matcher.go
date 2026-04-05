package lanzou

import (
	"encoding/json"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

// AjaxData 从页面提取的AJAX请求参数
type AjaxData struct {
	URL  string
	Type string // GET/POST
	Data map[string]string
}

// Matcher 从蓝奏云HTML页面提取各种数据
type Matcher struct{}

// MatchIframe 从HTML中提取iframe src
func MatchIframe(html string) string {
	re := regexp.MustCompile(`<iframe[^>]+src=["']([^"']+)["']`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

// ParsePwdAjax 从带密码的分享页面提取AJAX参数
func ParsePwdAjax(html, pwd string) (*AjaxData, error) {
	// 提取url和sign等参数
	urlRe := regexp.MustCompile(`url\s*:\s*["']([^"']+)["']`)
	signRe := regexp.MustCompile(`sign\s*:\s*["']([^"']*)["']`)

	urlMatch := urlRe.FindStringSubmatch(html)
	signMatch := signRe.FindStringSubmatch(html)

	ajaxURL := ""
	if len(urlMatch) > 1 {
		ajaxURL = urlMatch[1]
	}

	// 构建POST数据
	data := map[string]string{
		"action": "load2",
		"passwd": pwd,
	}

	if len(signMatch) > 1 && signMatch[1] != "" {
		data["sign"] = signMatch[1]
	}

	if ajaxURL == "" {
		return nil, fmt.Errorf("无法提取AJAX URL")
	}

	return &AjaxData{
		URL:  ajaxURL,
		Type: "POST",
		Data: data,
	}, nil
}

// ParseFolderAjax 从文件夹分享页面提取AJAX参数
func ParseFolderAjax(html string) (*AjaxData, error) {
	// 提取AJAX URL
	urlRe := regexp.MustCompile(`url\s*:\s*["']([^"']+)["']`)
	urlMatch := urlRe.FindStringSubmatch(html)

	// 提取pgs变量
	pgsRe := regexp.MustCompile(`var\s+pgs\s*=\s*(\d+)`)
	pgsMatch := pgsRe.FindStringSubmatch(html)

	ajaxURL := ""
	if len(urlMatch) > 1 {
		ajaxURL = urlMatch[1]
	}

	pgs := "1"
	if len(pgsMatch) > 1 {
		pgs = pgsMatch[1]
	}

	if ajaxURL == "" {
		return nil, fmt.Errorf("无法提取文件夹AJAX URL")
	}

	data := map[string]string{
		"pg": pgs,
	}

	// 提取其他参数
	veiRe := regexp.MustCompile(`vei\s*:\s*["']([^"']*)["']`)
	if m := veiRe.FindStringSubmatch(html); len(m) > 1 && m[1] != "" {
		data["vei"] = m[1]
	}

	return &AjaxData{
		URL:  ajaxURL,
		Type: "POST",
		Data: data,
	}, nil
}

// ParseFileMoreAjax 从"更多文件"页面提取AJAX参数
func ParseFileMoreAjax(html string) (*AjaxData, error) {
	urlRe := regexp.MustCompile(`url\s*:\s*["']([^"']+)["']`)
	urlMatch := urlRe.FindStringSubmatch(html)

	if len(urlMatch) < 2 || urlMatch[1] == "" {
		return nil, fmt.Errorf("无法提取AJAX URL")
	}

	data := map[string]string{
		"pg":        "1",
		"folder_id": "",
	}

	return &AjaxData{
		URL:  urlMatch[1],
		Type: "POST",
		Data: data,
	}, nil
}

// BuildQueryString 将map构建为URL编码的查询字符串
func BuildQueryString(data map[string]string) string {
	values := url.Values{}
	for k, v := range data {
		values.Set(k, v)
	}
	return values.Encode()
}

// ExtractShareFileInfo 从分享页面HTML提取文件基本信息
func ExtractShareFileInfo(html, shareURL, pwd string) (*ShareFileInfo, error) {
	name := ExtractShareTitle(html)
	size := ExtractShareFileSize(html)
	timeStr := ExtractShareFileTime(html)

	return &ShareFileInfo{
		Name: name,
		Size: size,
		Time: timeStr,
		URL:  shareURL,
		Pwd:  pwd,
	}, nil
}

// ExtractShareTitle 从分享页面提取文件名
func ExtractShareTitle(html string) string {
	// 优先从title提取
	re := regexp.MustCompile(`<title>(.+?)</title>`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(strings.TrimSuffix(match[1], " - 蓝奏云"))
	}

	// 从span.b提取
	re = regexp.MustCompile(`<span class="b">([^<]+)</span>`)
	match = re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	return ""
}

// ExtractShareFileSize 从分享页面提取文件大小
func ExtractShareFileSize(html string) string {
	// 从meta description提取: "文件大小：100MB|..."
	re := regexp.MustCompile(`content=["'][^"']*文件大小[：:]([^|]+)\|`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	return "未知"
}

// ExtractShareFileTime 从分享页面提取上传时间
func ExtractShareFileTime(html string) string {
	// 从 n_file_infos 提取
	re := regexp.MustCompile(`n_file_infos[^>]*>([^<]+)<`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	// 尝试"上传时间：xxx<br>"
	re = regexp.MustCompile(`上传时间[：：</span>]*([^<]+)<br`)
	match = re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSpace(match[1])
	}

	return ""
}

// DetectShareType 检测分享页面类型
type SharePageType int

const (
	ShareTypeFile SharePageType = iota // 无密码文件
	ShareTypePwdFile                   // 有密码文件
	ShareTypeFolder                    // 无密码文件夹
	ShareTypePwdFolder                 // 有密码文件夹
	ShareTypeUnknown                   // 未知类型
)

// DetectSharePageType 检测分享页面类型
func DetectSharePageType(html string) SharePageType {
	hasIframe := strings.Contains(html, "<iframe")
	hasPasswdDiv := strings.Contains(html, "passwddiv")
	hasFilemore := strings.Contains(html, "filemore")
	hasPwdload := strings.Contains(html, "pwdload")

	if hasPasswdDiv && !hasIframe {
		return ShareTypePwdFile
	}
	if hasIframe {
		return ShareTypeFile
	}
	if hasPwdload {
		return ShareTypePwdFolder
	}
	if hasFilemore {
		return ShareTypeFolder
	}
	return ShareTypeUnknown
}

// ParseShareFolderResponse 解析文件夹分享AJAX响应
func ParseShareFolderResponse(jsonData []byte) ([]ShareFileInfo, error) {
	var resp struct {
		Zt   int             `json:"zt"`
		Info string          `json:"info"`
		Text json.RawMessage `json:"text"`
	}

	if err := json.Unmarshal(jsonData, &resp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.Zt != 1 {
		return nil, fmt.Errorf("请求失败: %s", resp.Info)
	}

	// 蓝奏云API返回的文件夹文件列表使用 name_all 字段
	var files []struct {
		ID      string `json:"id"`
		NameAll string `json:"name_all"`
		Name    string `json:"name"`
		Size    string `json:"size"`
		Time    string `json:"time"`
		Icon    string `json:"icon"`
	}
	if err := json.Unmarshal(resp.Text, &files); err != nil {
		return nil, fmt.Errorf("解析文件列表失败: %w", err)
	}

	result := make([]ShareFileInfo, 0, len(files))
	for _, f := range files {
		name := f.NameAll
		if name == "" {
			name = f.Name
		}
		result = append(result, ShareFileInfo{
			Name: name,
			Size: f.Size,
			Time: f.Time,
			URL:  "",
		})
	}

	return result, nil
}
