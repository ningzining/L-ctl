package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/util/templateutil"
	"github.com/spf13/cobra"
	"os/exec"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新模板库",
	Long:  `更新模板库`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateTemplate()
	},
}

func init() {
	templateCmd.AddCommand(updateCmd)
}

func updateTemplate() error {
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
