package templateutil

import (
	"errors"
	"fmt"
	"github.com/ningzining/L-ctl/logic/version"
	"net/url"
	"os"
	"runtime"
	"text/template"
)

const (
	TemplateGitUrl  = "https://github.com/ningzining/L-ctl-template.git"
	TemplateRepoUrl = "https://raw.githubusercontent.com/ningzining/L-ctl-template/main/repo/repo.tpl"
	LocalRepoUrl    = "repo/repo.tpl"
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

// GetLocalRepoTemplate 获取本地repo.tpl模板文件的路径
func GetLocalRepoTemplate() (string, error) {
	dir, err := GenerateTemplateDir()
	if err != nil {
		return "", err
	}
	result, err := url.JoinPath(dir, LocalRepoUrl)
	if err != nil {
		return "", err
	}
	return result, nil
}

// SaveTemplateByLocal 渲染数据到指定的模板，并保存
// templatePath: 模板路径
// filePath: 保存的路径
// data: 数据源
func SaveTemplateByLocal(templatePath string, filePath string, data interface{}) error {
	templateFiles, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	err = templateFiles.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}

// SaveTemplateByData 渲染数据到指定的模板，并保存
// templateData: 模板字节
// filePath: 保存的路径
// data: 数据源
func SaveTemplateByData(templateData []byte, filePath string, data interface{}) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 根据字节数组创建模板
	files, err := template.New("temp").Parse(string(templateData))
	if err != nil {
		return err
	}
	// 渲染数据到模板
	err = files.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}
