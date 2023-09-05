package swag

import (
	"github.com/ningzining/L-ctl/logic"
	"github.com/spf13/cobra"
)

var file string
var projectId string

var Cmd = &cobra.Command{
	Use:   "swag",
	Short: "导入swagger文档到指定apifox",
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := logic.SwagGenerateArg{
			File:      file,
			ProjectId: projectId,
		}
		return logic.NewSwag().Upload(arg)
	},
}

func init() {
	Cmd.Flags().StringVarP(&file, "file", "f", "", "指定swagger文件的位置")
	Cmd.Flags().StringVarP(&projectId, "projectId", "p", "", "指定apifox当中的项目")
	if err := Cmd.MarkFlagRequired("file"); err != nil {
		return
	}
	if err := Cmd.MarkFlagRequired("projectId"); err != nil {
		return
	}
}
