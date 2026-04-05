package service

import (
	"github.com/wenaixi/wenxi-cloud/backend/internal/pkg/lanzou"
)

type ShareParseService struct {
	client *lanzou.Client
}

func NewShareParseService(client *lanzou.Client) *ShareParseService {
	return &ShareParseService{client: client}
}

// ParseShareLink 解析分享链接
func (s *ShareParseService) ParseShareLink(shareURL, pwd string) (interface{}, error) {
	return s.client.ParseShareURL(shareURL, pwd)
}

// BatchParseShareLinks 批量解析分享链接
func (s *ShareParseService) BatchParseShareLinks(urls []string, pwd string) ([]interface{}, []string) {
	results := make([]interface{}, 0, len(urls))
	errors := make([]string, 0)

	for _, url := range urls {
		if url == "" {
			continue
		}
		result, err := s.ParseShareLink(url, pwd)
		if err != nil {
			errors = append(errors, err.Error())
		} else {
			results = append(results, result)
		}
	}
	return results, errors
}

// GetShareDownloadLink 获取分享文件的下载直链
func (s *ShareParseService) GetShareDownloadLink(shareURL, pwd string) (string, error) {
	return s.client.GetShareDownloadURL(shareURL, pwd)
}

// ValidateShareLink 验证分享链接是否有效
func (s *ShareParseService) ValidateShareLink(shareURL string) bool {
	return lanzou.ValidateShareURL(shareURL)
}
