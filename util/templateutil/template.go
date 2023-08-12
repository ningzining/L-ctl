package templateutil

import (
	"github.com/ningzining/L-ctl/logic/version"
	"os"
	"path/filepath"
	"text/template"
)

const (
	TemplateGitUrl     = "https://github.com/ningzining/L-ctl-template.git"
	TemplateRepoUrl    = "https://raw.githubusercontent.com/ningzining/L-ctl-template/main/repo/repo.tpl"
	LocalRepoUrl       = "repo/repo.tpl"
	DefaultTemplateDir = ".L-ctl"
)

// GenerateTemplateDir 生成默认模板文件的路径
func GenerateTemplateDir() (string, error) {
	dirPrefix, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	resultDir := filepath.Join(dirPrefix, DefaultTemplateDir, version.BuildVersion)
	return resultDir, nil
}

// GetRepoTemplatePath 获取本地repo.tpl模板文件的路径
func GetRepoTemplatePath() (string, error) {
	dir, err := GenerateTemplateDir()
	if err != nil {
		return "", err
	}
	result := filepath.Join(dir, LocalRepoUrl)
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
