package logic

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/util/pathutil"
	"github.com/ningzining/L-ctl/util/templateutil"
	"os/exec"
	"path/filepath"
)

type Template struct{}

func NewTemplate() *Template {
	return &Template{}
}

// Init 初始化模板文件
func (t *Template) Init() error {
	// 获取模板文件所存在的目录
	templateDir, err := templateutil.GenerateTemplateDir()
	if err != nil {
		return err
	}
	exist, err := pathutil.Exist(templateDir)
	if err != nil {
		return err
	}
	if !exist {
		// 不存在目录则创建目标目录
		err := pathutil.Mkdir(templateDir)
		if err != nil {
			return err
		}
	}
	// 判断模板文件是否被初始化过
	localDir := filepath.Join(templateDir, ".git")
	if b, err := pathutil.Exist(localDir); b || err != nil {
		return errors.New(fmt.Sprintf("目标文件夹不为空,模板初始化失败,请检查目标文件夹: %s\n", templateDir))
	}
	// clone模板库到本地
	command := exec.Command("git", "clone", templateutil.TemplateGitUrl, templateDir)
	output, err := command.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("模板初始化失败: %s", err.Error()))
	}
	fmt.Printf("%s\n", string(output))
	color.Green("模板库初始化成功: %s\n", templateDir)
	return nil
}

// Update 更新模板文件
func (t *Template) Update() error {
	templateDir, err := templateutil.GenerateTemplateDir()
	if err != nil {
		return err
	}
	cmd := exec.Command("git", "pull")
	cmd.Dir = templateDir
	bytes, err := cmd.Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(bytes))
	color.Green("模板库更新成功")
	return nil
}
