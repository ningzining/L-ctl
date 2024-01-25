package swag

import (
	"github.com/fatih/color"
	"github.com/ningzining/lctl/logic/swag"
	"github.com/spf13/cobra"
)

var (
	file      string
	projectId string
)

var CmdSwag = &cobra.Command{
	Use:   "swag",
	Short: "导入swagger文档到指定apifox",
	Run: func(cmd *cobra.Command, args []string) {
		if err := swag.NewSwag(file, projectId).Upload(); err != nil {
			color.Red("swagger导入apifox失败: %s\n", err.Error())
			return
		}
	},
}

func init() {
	CmdSwag.Flags().StringVarP(&file, "file", "f", "", "指定swagger.json文件的位置")
	CmdSwag.Flags().StringVarP(&projectId, "projectId", "p", "", "指定apifox当中的项目ID")
	if err := CmdSwag.MarkFlagRequired("file"); err != nil {
		return
	}
	if err := CmdSwag.MarkFlagRequired("projectId"); err != nil {
		return
	}
}
