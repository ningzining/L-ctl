package template

import (
	"github.com/fatih/color"
	"github.com/ningzining/lctl/logic/template"
	"github.com/spf13/cobra"
)

var CmdTemplate = &cobra.Command{
	Use:   "template",
	Short: "操作模板库",
	Long:  "操作模板库，提供了初始化和更新的操作",
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化模板库",
	Long:  `初始化模板库`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := template.NewTemplate().Init(); err != nil {
			color.Red("模板初始化失败,%s", err.Error())
			return err
		}
		return nil
	},
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新模板库",
	Long:  `更新模板库`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return template.NewTemplate().Update()
	},
}

func init() {
	CmdTemplate.AddCommand(initCmd)
	CmdTemplate.AddCommand(updateCmd)
}
