package swag

import (
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/logic/swag"
	"github.com/spf13/cobra"
)

var (
	file      string
	projectId string
)

var Cmd = &cobra.Command{
	Use:   "swag",
	Short: "导入swagger文档到指定apifox",
	Run: func(cmd *cobra.Command, args []string) {
		if err := swag.NewSwag(file, projectId).Upload(); err != nil {
			color.Red("swagger导入apifox失败: %s\n", err.Error())
			return
		}
		color.Green("swagger导入apifox成功\n")
	},
}

func init() {
	Cmd.Flags().StringVarP(&file, "file", "f", "", "指定swagger.json文件的位置")
	Cmd.Flags().StringVarP(&projectId, "projectId", "p", "", "指定apifox当中的项目ID")
	if err := Cmd.MarkFlagRequired("file"); err != nil {
		return
	}
	if err := Cmd.MarkFlagRequired("projectId"); err != nil {
		return
	}
}
