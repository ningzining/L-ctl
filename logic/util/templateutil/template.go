package templateutil

import (
	"github.com/ningzining/L-ctl/logic/version"
	"net/url"
	"runtime"
)

func GenerateTemplateDir() (string, error) {
	var dirPrefix string
	tempDir := ".L-ctl"
	switch runtime.GOOS {
	case "windows":
		dirPrefix = "C:/Users/Admin"
	case "darwin":
		dirPrefix = "~/"
	}
	resultDir, err := url.JoinPath(dirPrefix, tempDir, version.BuildVersion)
	if err != nil {
		return "", err
	}
	return resultDir, nil
}
