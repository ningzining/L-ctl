package template

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ningzining/lctl/util/gitutil"
	"github.com/ningzining/lctl/util/pathutil"
	"github.com/ningzining/lctl/util/templateutil"
)

// Init 初始化模板文件
func (t *Template) Init() error {
	// 获取模板文件所存在的目录
	templateDir, err := templateutil.GetTemplateDir()
	if err != nil {
		return err
	}

	// 不存在目录则创建目标目录
	if err := pathutil.MkdirIfNotExist(templateDir); err != nil {
		return err
	}

	// 判断模板文件是否被初始化过
	if b, err := pathutil.Exist(filepath.Join(templateDir, ".git")); b || err != nil {
		return errors.New(fmt.Sprintf("目标文件夹不为空,模板初始化失败,请检查目标文件夹: %s\n", templateDir))
	}

	// clone模板库到本地
	output, err := gitutil.CloneTemplate(templateDir)
	if err != nil {
		return errors.New(fmt.Sprintf("模板初始化失败: %s", err.Error()))
	}

	fmt.Println(string(output))
	color.Green("模板库初始化成功: %s", templateDir)
	return nil
}
