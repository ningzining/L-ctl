package templateutil

import (
	"bytes"
	"errors"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const (
	TemplateGitUrl     = "https://github.com/ningzining/L-ctl-template.git"
	LocalModelUrl      = "model/model.tpl"
	LocalRepoUrl       = "repo/repo.tpl"
	DefaultTemplateDir = ".L-ctl"
	Template           = "template"
	GoFileSuffix       = ".go"
)

// GenerateTemplateDir 生成默认模板文件的路径
func GenerateTemplateDir() (string, error) {
	dirPrefix, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	resultDir := filepath.Join(dirPrefix, DefaultTemplateDir, Template)
	return resultDir, nil
}

// GetTemplateDir 获取默认模板文件的路径
func GetTemplateDir() (string, error) {
	dirPrefix, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	resultDir := filepath.Join(dirPrefix, DefaultTemplateDir, Template)
	return resultDir, nil
}

// GetConfigFilePath 获取默认配置文件的路径
func GetConfigFilePath() (string, error) {
	dirPrefix, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dirPrefix, DefaultTemplateDir), nil
}

// Create 创建模板文件
func Create(filePath string, data map[string]any, types string) (err error) {
	var templatePath string
	switch types {
	case LocalModelUrl:
		// 获取本地模板文件的路径
		templatePath, err = GetModelTemplatePath()
		if err != nil {
			return err
		}
	case LocalRepoUrl:
		templatePath, err = GetRepoTemplatePath()
		if err != nil {
			return err
		}
	default:
		return errors.New("模板类型不存在")
	}
	// 通过本地文件保存模板
	if err := SaveTemplateByLocal(templatePath, filePath, data); err != nil {
		return err
	}
	return nil
}

// CreateRepoFile 创建repo模板文件
func CreateRepoFile(filePath string, data map[string]any) (err error) {
	// 获取repo模板文件路径
	templatePath, err := GetRepoTemplatePath()
	if err != nil {
		return err
	}
	// 创建模板文件
	if err := CreateTemplateFile(templatePath, filePath, data); err != nil {
		return err
	}
	return nil
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

// GetModelTemplatePath 获取本地model.tpl模板文件的路径
func GetModelTemplatePath() (string, error) {
	dir, err := GenerateTemplateDir()
	if err != nil {
		return "", err
	}
	result := filepath.Join(dir, LocalModelUrl)
	return result, nil
}

// SaveTemplateByLocal 渲染数据到指定的模板，并保存
// templatePath: 模板路径
// filePath: 保存的路径
// data: 数据源
func SaveTemplateByLocal(templatePath string, filePath string, data map[string]any) error {
	templateFile, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	// 渲染数据到模板
	buf := new(bytes.Buffer)
	err = templateFile.Execute(buf, data)
	if err != nil {
		return err
	}

	// 创建文件
	if err = createFile(filePath, buf.Bytes()); err != nil {
		return err
	}
	return nil
}

// CreateTemplateFile 创建模板文件
func CreateTemplateFile(templatePath string, filePath string, data map[string]any) error {
	templateFile, err := template.ParseFiles(templatePath)
	if err != nil {
		return err
	}

	// 渲染数据到模板
	buf := new(bytes.Buffer)
	err = templateFile.Execute(buf, data)
	if err != nil {
		return err
	}

	// 创建文件
	return createFile(filePath, buf.Bytes())
}

// 创建文件
// filePath: 文件路径
// dadaBytes: 字节数组
func createFile(filePath string, dataByes []byte) error {
	buf := new(bytes.Buffer)
	// 如果是创建go文件，则进行一次goFormat格式化
	ext := path.Ext(filePath)
	if ext == GoFileSuffix {
		formatOutput, err := format.Source(dataByes)
		if err != nil {
			buf.Write(dataByes)
		} else {
			buf.Write(formatOutput)
		}
	} else {
		buf.Write(dataByes)
	}

	return os.WriteFile(filePath, buf.Bytes(), os.ModePerm)
}

// MergeMap 合并所有的map集合
func MergeMap(source ...map[string]any) map[string]any {
	res := make(map[string]any)
	for _, t := range source {
		for k, v := range t {
			res[k] = v
		}
	}
	return res
}
