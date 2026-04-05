package lanzou

import (
	"crypto/rand"
	"math/big"
	"strings"
)

// SupportedExtensions 蓝奏云支持的文件扩展名列表
var SupportedExtensions = []string{
	"doc", "docx", "zip", "rar", "apk", "ipa", "txt", "dmg", "exe", "7z", "e", "z",
	"ct", "ke", "cetrainer", "db", "tar", "pdf", "w3x", "epub", "mobi", "azw", "azw3",
	"osk", "osz", "xpa", "cpk", "lua", "jar", "ppt", "pptx", "xls", "xlsx", "mp3",
	"iso", "img", "gho", "ttf", "ttc", "txf", "dwg", "bat", "imazingapp", "dll", "crx",
	"xapk", "conf", "deb", "rp", "rpm", "rplib", "mobileconfig", "appimage", "lolgezi",
	"flac", "cad", "hwt", "accdb", "ce", "xmind", "enc", "bds", "bdi", "ssf", "it",
	"pkg", "cfg",
}

// SafeSuffixes 安全后缀列表，用于文件名混淆
var SafeSuffixes = []string{
	"ct", "ke", "w3x", "mobi", "azw", "azw3", "osk", "osz", "xpa", "cpk",
	"lua", "gho", "ttc", "txf", "bat", "imazingapp", "xapk", "conf", "rp",
	"rplib", "mobileconfig", "appimage", "lolgezi", "cad", "hwt", "ce",
	"xmind", "bds", "bdi", "ssf", "it", "pkg", "cfg",
}

// IsSupportedExtension 检查文件扩展名是否被蓝奏云支持
func IsSupportedExtension(ext string) bool {
	ext = strings.ToLower(strings.TrimPrefix(ext, "."))
	for _, supported := range SupportedExtensions {
		if ext == supported {
			return true
		}
	}
	return false
}

// CreateSpecificName 为不支持的文件名添加随机安全后缀
// 例如: program.exe -> program.exe.lua.w3x
func CreateSpecificName(fileName string) string {
	suffix1 := randomSuffix()
	suffix2 := randomSuffix()
	return fileName + "." + suffix1 + "." + suffix2
}

// randomSuffix 从安全后缀列表中随机选择一个
func randomSuffix() string {
	idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(SafeSuffixes))))
	return SafeSuffixes[idx.Int64()]
}

// GetFileExtension 获取文件扩展名（不带点）
func GetFileExtension(fileName string) string {
	lastDot := strings.LastIndex(fileName, ".")
	if lastDot == -1 || lastDot == len(fileName)-1 {
		return ""
	}
	return strings.ToLower(fileName[lastDot+1:])
}
