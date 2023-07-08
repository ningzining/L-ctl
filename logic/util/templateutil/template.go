package templateutil

import (
	"errors"
	"fmt"
	"github.com/ningzining/L-ctl/logic/version"
	"net/url"
	"runtime"
)

// GenerateTemplateDir 生成默认模板文件的路径
func GenerateTemplateDir() (string, error) {
	var dirPrefix string
	tempDir := ".L-ctl"
	switch runtime.GOOS {
	case "windows":
		dirPrefix = "C:/Users/Admin"
	case "darwin":
		dirPrefix = "~/"
	default:
		return "", errors.New(fmt.Sprintf("目前系统不支持该操作系统: %s\n", runtime.GOOS))
	}
	resultDir, err := url.JoinPath(dirPrefix, tempDir, version.BuildVersion)
	if err != nil {
		return "", err
	}
	return resultDir, nil
}
