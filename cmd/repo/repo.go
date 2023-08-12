package repo

import (
	"github.com/ningzining/L-ctl/logic"
	"github.com/spf13/cobra"
)

var dir string
var table string
var style string

var Cmd = &cobra.Command{
	Use:   "repo",
	Short: "生成仓储层的repo文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		return logic.NewRepo().Generate(dir, table, style)
	},
}

func init() {
	Cmd.Flags().StringVarP(&dir, "dir", "d", "", "指定文件生成的目录")
	Cmd.Flags().StringVarP(&table, "table", "t", "", "指定数据库的表名")
	Cmd.Flags().StringVarP(&style, "style", "s", "", "指定生成的文件格式,默认为下划线格式")
	if err := Cmd.MarkFlagRequired("dir"); err != nil {
		return
	}
	if err := Cmd.MarkFlagRequired("table"); err != nil {
		return
	}
}
