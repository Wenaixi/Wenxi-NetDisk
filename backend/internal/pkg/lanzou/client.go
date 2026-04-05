package lanzou

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Client 蓝奏云API客户端
type Client struct {
	baseURL    string        // 基础URL，如 https://pc.woozooo.com
	httpClient *http.Client  // HTTP客户端
	cookie     string        // Cookie字符串
	userAgent  string        // User-Agent
}

// NewClient 创建蓝奏云客户端
func NewClient(cookie string) *Client {
	jar, _ := cookiejar.New(nil)
	return &Client{
		baseURL: "https://pc.woozooo.com",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
		},
		cookie:    cookie,
		userAgent: "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}
}

// SetCookie 设置Cookie
func (c *Client) SetCookie(cookie string) {
	c.cookie = cookie
}

// SetBaseURL 设置基础URL（用于测试）
func (c *Client) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// buildRequest 构建请求
func (c *Client) buildRequest(method, path string, body url.Values) (*http.Request, error) {
	reqURL := c.baseURL + "/" + path
	req, err := http.NewRequest(method, reqURL, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Referer", c.baseURL+"/")
	req.Header.Set("Origin", c.baseURL)

	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}

	return req, nil
}

// doRequest 发送请求
func (c *Client) doRequest(method, path string, body url.Values) ([]byte, error) {
	req, err := c.buildRequest(method, path, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// postForm 发送POST表单请求
func (c *Client) postForm(path string, body url.Values) ([]byte, error) {
	return c.doRequest("POST", path, body)
}

// Task5 列出文件列表
// folderId: 文件夹ID，-1表示根目录
// pg: 页码，从1开始
func (c *Client) Task5(folderId int, pg int) (*Task5Response, error) {
	body := url.Values{
		"task":      {"5"},
		"folder_id": {strconv.Itoa(folderId)},
		"pg":        {strconv.Itoa(pg)},
		"vei":       {"e11ad"}, // 固定值
	}

	data, err := c.postForm("doupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task5Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task47 列出文件夹列表
func (c *Client) Task47(folderId int) (*Task47Response, error) {
	body := url.Values{
		"task":      {"47"},
		"folder_id": {strconv.Itoa(folderId)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task47Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task2 创建文件夹
// parentId: 父文件夹ID，-1表示根目录
// name: 文件夹名称
func (c *Client) Task2(parentId int, name string) (*Task2Response, error) {
	// 蓝奏云会把空格和括号替换掉
	name = regexp.MustCompile(`[ ()]`).ReplaceAllString(name, "_")

	body := url.Values{
		"task":               {"2"},
		"parent_id":          {strconv.Itoa(parentId)},
		"folder_name":        {name},
		"folder_description": {""},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task2Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task6 删除文件
func (c *Client) Task6(fileId int) (*DeleteResponse, error) {
	body := url.Values{
		"task":    {"6"},
		"file_id": {strconv.Itoa(fileId)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &DeleteResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task46 删除文件夹
func (c *Client) Task46(folderId int) (*DeleteResponse, error) {
	body := url.Values{
		"task":      {"46"},
		"folder_id": {strconv.Itoa(folderId)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &DeleteResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task14 重命名文件
func (c *Client) Task14(fileId int, name string) (*RenameResponse, error) {
	body := url.Values{
		"task":    {"14"},
		"file_id": {strconv.Itoa(fileId)},
		"name":    {name},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &RenameResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task15 移动文件
func (c *Client) Task15(fileId, targetFolderId int) (*MoveResponse, error) {
	body := url.Values{
		"task":      {"15"},
		"file_id":   {strconv.Itoa(fileId)},
		"folder_id": {strconv.Itoa(targetFolderId)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &MoveResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task48 移动文件夹
func (c *Client) Task48(folderId, targetFolderId int) (*MoveResponse, error) {
	body := url.Values{
		"task":      {"48"},
		"folder_id": {strconv.Itoa(folderId)},
		"target_id": {strconv.Itoa(targetFolderId)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &MoveResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task48FolderRename 重命名文件夹 (task=4)
func (c *Client) Task48FolderRename(folderId int, name string) (*Task4Response, error) {
	name = regexp.MustCompile(`[ ()]`).ReplaceAllString(name, "_")

	body := url.Values{
		"task":      {"4"},
		"folder_id": {strconv.Itoa(folderId)},
		"folder_name": {name},
		"folder_description": {""},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task4Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task22 文件详情
func (c *Client) Task22(fileId int) (map[string]interface{}, error) {
	body := url.Values{
		"task":    {"22"},
		"file_id": {strconv.Itoa(fileId)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Task18 文件夹详情
func (c *Client) Task18(folderId int) (map[string]interface{}, error) {
	body := url.Values{
		"task":      {"18"},
		"folder_id": {strconv.Itoa(folderId)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// Task39 创建分享链接
// minutes: 有效期(分钟)，0表示永久
func (c *Client) Task39(fileId int, minutes int) (*ShareResponse, error) {
	body := url.Values{
		"task":     {"39"},
		"file_id":  {strconv.Itoa(fileId)},
		"onetime":  {strconv.Itoa(minutes)},
	}

	data, err := c.postForm("douupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &ShareResponse{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// GetDownloadURL 获取文件下载直链
// shareURL: 蓝奏云分享链接，如 https://wws.lanzous.com/xxxxx
// pwd: 密码（如果有）
func (c *Client) GetDownloadURL(shareURL, pwd string) (string, string, error) {
	// 解析分享页面
	resp, err := c.httpClient.Get(shareURL)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	htmlStr := string(html)

	// 检查是否有密码
	if strings.Contains(htmlStr, "passwddiv") && pwd == "" {
		return "", "", fmt.Errorf("share requires password")
	}

	// 提取文件名
	name := extractName(htmlStr)

	// 提取下载直链
	downURL := extractDownloadURL(htmlStr, shareURL, pwd)

	return name, downURL, nil
}

// extractName 从HTML提取文件名
func extractName(html string) string {
	re := regexp.MustCompile(`<title>(.+?)</title>`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		return strings.TrimSuffix(match[1], " - 蓝奏云")
	}
	return ""
}

// extractDownloadURL 从HTML提取下载链接
func extractDownloadURL(html, referer, pwd string) string {
	// 简化实现，实际需要解析iframe和ajax
	re := regexp.MustCompile(`iframe[^>]+src=["']([^"']+)["']`)
	match := re.FindStringSubmatch(html)
	if len(match) > 1 {
		iframeURL := match[1]
		if !strings.HasPrefix(iframeURL, "http") {
			iframeURL = "https://wws.lanzous.com" + iframeURL
		}
		return iframeURL
	}
	return ""
}

// doGet 发送GET请求
func (c *Client) doGet(path string) ([]byte, error) {
	reqURL := c.baseURL + "/" + path
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", c.baseURL+"/")

	if c.cookie != "" {
		req.Header.Set("Cookie", c.cookie)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// Profile 获取用户个人信息
// 从 mydisk.php 页面解析 iframe referer
// 从 mydisk.php?item=profile&action=mypower 页面解析用户信息
func (c *Client) Profile() (*ProfileInfo, error) {
	// 并发请求两个页面
	type result struct {
		data []byte
		err  error
	}
	ch := make(chan result, 2)

	go func() {
		data, err := c.doGet("mydisk.php")
		ch <- result{data, err}
	}()

	go func() {
		data, err := c.doGet("mydisk.php?item=profile&action=mypower")
		ch <- result{data, err}
	}()

	r1 := <-ch
	r2 := <-ch

	if r1.err != nil {
		return nil, fmt.Errorf("failed to fetch mydisk.php: %w", r1.err)
	}
	if r2.err != nil {
		return nil, fmt.Errorf("failed to fetch profile page: %w", r2.err)
	}

	profile := &ProfileInfo{}

	// 解析 profile 页面
	c.parseProfilePage(string(r2.data), profile)

	// 解析主页面获取 referer
	profile.Referer = extractRefererFromMainPage(string(r1.data), c.baseURL)

	return profile, nil
}

// parseProfilePage 从 profile HTML 页面解析用户信息
func (c *Client) parseProfilePage(html string, profile *ProfileInfo) {
	// 解析个性域名 - <span id="domaindiynow">xxx</span>
	domainRe := regexp.MustCompile(`<span[^>]*id=["']domaindiynow["'][^>]*>(.*?)</span>`)
	if m := domainRe.FindStringSubmatch(html); len(m) > 1 {
		profile.Domain = strings.TrimSpace(m[1])
	}

	// 解析最近登录时间 - 使用[\s\S]*?代替.*?来匹配跨行
	lastLoginRe := regexp.MustCompile(`最近登录时间[\s\S]*?<span[^>]*class=["'][^"']*mf2[^"']*["'][^>]*>(.*?)</span>`)
	if m := lastLoginRe.FindStringSubmatch(html); len(m) > 1 {
		profile.LastLogin = strings.TrimSpace(m[1])
	}

	// 解析允许上传类型 - 跨行匹配
	supportRe := regexp.MustCompile(`允许上传类型[\s\S]*?<span[^>]*class=["'][^"']*mf2[^"']*["'][^>]*>(.*?)</span>`)
	if m := supportRe.FindStringSubmatch(html); len(m) > 1 {
		raw := m[1]
		raw = strings.ReplaceAll(raw, "<br>", ",")
		raw = strings.ReplaceAll(raw, "<br/>", ",")
		raw = strings.ReplaceAll(raw, "<br />", ",")
		parts := strings.Split(raw, ",")
		for _, p := range parts {
			p = strings.TrimSpace(p)
			if p != "" {
				profile.SupportList = append(profile.SupportList, p)
			}
		}
	}

	// 解析单个文件大小 - 跨行匹配
	maxSizeRe := regexp.MustCompile(`单个文件大小[\s\S]*?<font[^>]*>(.*?)</font>`)
	if m := maxSizeRe.FindStringSubmatch(html); len(m) > 1 {
		profile.MaxSize = strings.TrimSpace(m[1])
	}

	// 解析安全验证 - <span id="phone_id">xxx</span>
	verificationRe := regexp.MustCompile(`<span[^>]*id=["']phone_id["'][^>]*>(.*?)</span>`)
	if m := verificationRe.FindStringSubmatch(html); len(m) > 1 {
		profile.Verification = strings.TrimSpace(m[1])
	}
}

// extractRefererFromMainPage 从主页面提取 iframe src 作为 referer
func extractRefererFromMainPage(html, baseURL string) string {
	// 匹配 iframe src
	re := regexp.MustCompile(`<iframe[^>]+src=["']([^"']+)["']`)
	m := re.FindStringSubmatch(html)
	if len(m) > 1 {
		src := m[1]
		if !strings.HasPrefix(src, "http") {
			return baseURL + src
		}
		return src
	}
	return ""
}

// Ping 检测连接
func (c *Client) Ping() error {
	body := url.Values{"task": {"1"}}
	_, err := c.postForm("douload.php", body)
	return err
}

// Task23 设置文件访问密码
// fileId: 文件ID
// shows: 访问模式 (1=公开, 2=密码访问)
// shownames: 访问密码
func (c *Client) Task23(fileId int, shows int, shownames string) (*Task23Response, error) {
	body := url.Values{
		"task":       {"23"},
		"file_id":    {strconv.Itoa(fileId)},
		"shows":      {strconv.Itoa(shows)},
		"shownames":  {shownames},
	}

	data, err := c.postForm("doupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task23Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task12 获取文件描述
func (c *Client) Task12(fileId int) (*Task12Response, error) {
	body := url.Values{
		"task":    {"12"},
		"file_id": {strconv.Itoa(fileId)},
	}

	data, err := c.postForm("doupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task12Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task11 设置文件描述
func (c *Client) Task11(fileId int, desc string) (*Task11Response, error) {
	body := url.Values{
		"task":    {"11"},
		"file_id": {strconv.Itoa(fileId)},
		"desc":    {desc},
	}

	data, err := c.postForm("doupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task11Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}

// Task16 设置文件夹访问密码
// folderId: 文件夹ID
// shows: 访问模式 (1=公开, 2=密码访问)
// shownames: 访问密码
func (c *Client) Task16(folderId int, shows int, shownames string) (*Task16Response, error) {
	body := url.Values{
		"task":       {"16"},
		"folder_id":  {strconv.Itoa(folderId)},
		"shows":      {strconv.Itoa(shows)},
		"shownames":  {shownames},
	}

	data, err := c.postForm("doupload.php", body)
	if err != nil {
		return nil, err
	}

	resp := &Task16Response{}
	if err := json.Unmarshal(data, resp); err != nil {
		resp.Zt = -1
		resp.Info = string(data)
	}
	return resp, nil
}
