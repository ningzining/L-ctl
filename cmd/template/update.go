package template

import (
	"github.com/ningzining/L-ctl/logic"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "更新模板库",
	Long:  `更新模板库`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return logic.NewTemplate().Update()
	},
}

func init() {
	Cmd.AddCommand(updateCmd)
}
