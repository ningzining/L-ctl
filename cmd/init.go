package cmd

import (
	"errors"
	"fmt"
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

const templateUrl = "https://github.com/ningzining/L-ctl-template.git"

func initTemplate() error {
	templateDir, err := templateutil.GenerateTemplateDir()
	if err != nil {
		return err
	}
	exist, err := pathutil.Exist(templateDir)
	if err != nil {
		return err
	}
	if !exist {
		err := pathutil.Mkdir(templateDir)
		if err != nil {
			return err
		}
	}
	localDir, err := url.JoinPath(templateDir, ".git")
	if err != nil {
		return err
	}
	if b, err := pathutil.Exist(localDir); b || err != nil {
		return errors.New(fmt.Sprintf("目标文件夹不为空,模板初始化失败,请检查目标文件夹: %s\n", templateDir))
	}

	command := exec.Command("git", "clone", templateUrl, templateDir)
	output, err := command.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("模板初始化失败: %s", err.Error()))
	}
	fmt.Printf("%s\n", string(output))
	fmt.Printf("templates are generated in : %s\n", templateDir)
	return nil
}
