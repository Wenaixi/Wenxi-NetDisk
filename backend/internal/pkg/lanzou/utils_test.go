package lanzou

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestIsSupportedExtension 测试扩展名校验
func TestIsSupportedExtension(t *testing.T) {
	tests := []struct {
		ext      string
		expected bool
	}{
		{"txt", true},
		{"zip", true},
		{"pdf", true},
		{"doc", true},
		{"exe", true},
		{"mp3", true},
		{"iso", true},
		{"xyz", false},
		{"unknown", false},
		{"", false},
		{".txt", true},
		{".XYZ", false},
		{"TXT", true},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			result := IsSupportedExtension(tt.ext)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGetFileExtension 测试扩展名提取
func TestGetFileExtension(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"test.txt", "txt"},
		{"archive.tar.gz", "gz"},
		{"noextension", ""},
		{".hidden", "hidden"},
		{"file.", ""},
		{"FILE.TXT", "txt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetFileExtension(tt.name)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCreateSpecificName 测试文件名混淆
func TestCreateSpecificName(t *testing.T) {
	name := CreateSpecificName("program.exe")

	// 应该添加两个随机后缀: program.exe.suffix1.suffix2
	parts := splitByDot(name)
	assert.GreaterOrEqual(t, len(parts), 3) // program + exe + 2 random suffixes

	// 检查原始文件名保留（按点分割后第一个是program）
	assert.Equal(t, "program", parts[0])

	// 检查后缀在安全列表中（最后两个应该是安全后缀）
	suffix1 := parts[len(parts)-2]
	suffix2 := parts[len(parts)-1]
	assert.Contains(t, SafeSuffixes, suffix1)
	assert.Contains(t, SafeSuffixes, suffix2)
}

// TestCreateSpecificName_MultipleDots 测试多扩展名文件
func TestCreateSpecificName_MultipleDots(t *testing.T) {
	name := CreateSpecificName("archive.tar.gz")
	parts := splitByDot(name)

	// archive.tar.gz -> archive + tar + gz + suffix1 + suffix2
	assert.GreaterOrEqual(t, len(parts), 5)

	// 最后两个应该是安全后缀
	suffix1 := parts[len(parts)-2]
	suffix2 := parts[len(parts)-1]
	assert.Contains(t, SafeSuffixes, suffix1)
	assert.Contains(t, SafeSuffixes, suffix2)
}

// TestCreateSpecificName_NoExtension 测试无扩展名文件
func TestCreateSpecificName_NoExtension(t *testing.T) {
	name := CreateSpecificName("README")
	parts := splitByDot(name)

	assert.Equal(t, "README", parts[0])
	assert.GreaterOrEqual(t, len(parts), 3)
}

// TestSafeSuffixes_ContainsExpected 测试安全后缀列表
func TestSafeSuffixes_ContainsExpected(t *testing.T) {
	expected := []string{"lua", "w3x", "mobi", "azw", "gho", "ttc", "bat", "conf", "pkg", "cfg"}
	for _, suffix := range expected {
		assert.Contains(t, SafeSuffixes, suffix, "安全后缀列表应包含 %s", suffix)
	}
}

// TestSupportedExtensions_ContainsExpected 测试支持扩展列表
func TestSupportedExtensions_ContainsExpected(t *testing.T) {
	expected := []string{"txt", "zip", "rar", "7z", "pdf", "doc", "docx", "exe", "mp3", "iso"}
	for _, ext := range expected {
		assert.Contains(t, SupportedExtensions, ext, "支持扩展列表应包含 %s", ext)
	}
}

// TestRandomSuffix_ReturnsValidSuffix 测试随机后缀有效性
func TestRandomSuffix_ReturnsValidSuffix(t *testing.T) {
	for i := 0; i < 20; i++ {
		suffix := randomSuffix()
		assert.Contains(t, SafeSuffixes, suffix, "随机后缀 %d 应该是安全列表中的有效值", i)
	}
}

// splitByDot 辅助函数：按点分割字符串
func splitByDot(s string) []string {
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '.' {
			result = append(result, s[start:i])
			start = i + 1
		}
	}
	result = append(result, s[start:])
	return result
}
