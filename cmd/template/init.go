package template

import (
	"github.com/ningzining/L-ctl/logic"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化模板库",
	Long:  `初始化模板库`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return logic.NewTemplate().Init()
	},
}

func init() {
	Cmd.AddCommand(initCmd)
}
