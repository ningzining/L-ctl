package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/logic/util/pathutil"
	"github.com/ningzining/L-ctl/logic/util/templateutil"
	"github.com/spf13/cobra"
	"net/url"
	"os/exec"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化模板库",
	Long:  `初始化模板库`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return initTemplate()
	},
}

func init() {
	templateCmd.AddCommand(initCmd)
}

func initTemplate() error {
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
	localDir, err := url.JoinPath(templateDir, ".git")
	if err != nil {
		return err
	}
	if b, err := pathutil.Exist(localDir); b || err != nil {
		return errors.New(fmt.Sprintf("目标文件夹不为空,模板初始化失败,请检查目标文件夹: %s\n", templateDir))
	}

	command := exec.Command("git", "clone", templateutil.TemplateGitUrl, templateDir)
	output, err := command.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("模板初始化失败: %s", err.Error()))
	}
	fmt.Printf("%s\n", string(output))
	color.Green("templates are generated at : %s\n", templateDir)
	return nil
}
