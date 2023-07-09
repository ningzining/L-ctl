package cmd

import (
	"github.com/ningzining/L-ctl/logic"
	"github.com/spf13/cobra"
)

var dir string
var table string

var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "生成仓储层的repo文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logic.NewRepo().Generate(dir, table)
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.Flags().StringVarP(&dir, "dir", "d", "", "指定文件生成的目录")
	repoCmd.Flags().StringVarP(&table, "table", "t", "", "指定数据库的表名(文件名)")
	if err := repoCmd.MarkFlagRequired("dir"); err != nil {
		return
	}
	if err := repoCmd.MarkFlagRequired("table"); err != nil {
		return
	}
}
