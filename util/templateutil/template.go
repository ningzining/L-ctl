package templateutil

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path"
	"path/filepath"
	"text/template"
)

const (
	TemplateGitUrl     = "https://github.com/ningzining/L-ctl-template.git"
	TemplateRepoUrl    = "https://raw.githubusercontent.com/ningzining/L-ctl-template/main/repo/repo.tpl"
	LocalRepoUrl       = "repo/repo.tpl"
	LocalModelUrl      = "model/model.tpl"
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

// SaveTemplateByData 渲染数据到指定的模板，并保存
// templateData: 模板字节
// filePath: 保存的路径
// data: 数据源
func SaveTemplateByData(templateData []byte, filePath string, data interface{}) error {
	// 根据字节数组创建模板
	templateFile, err := template.New("temp").Parse(string(templateData))
	if err != nil {
		return err
	}

	// 渲染数据到模板
	buf := new(bytes.Buffer)
	err = templateFile.Execute(buf, data)

	// 创建文件
	if err = createFile(filePath, buf.Bytes()); err != nil {
		return err
	}
	return nil
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
			return errors.New(fmt.Sprintf("go文件格式化异常: %s\n", err))
		}
		buf.Write(formatOutput)
	} else {
		buf.Write(dataByes)
	}
	if err := os.WriteFile(filePath, buf.Bytes(), os.ModePerm); err != nil {
		return errors.New(fmt.Sprintf("文件创建失败: %s", err))
	}
	return nil
}
