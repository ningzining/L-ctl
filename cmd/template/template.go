package template

import (
	"github.com/ningzining/L-ctl/logic/template"
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
		return template.NewTemplate().Init()
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
