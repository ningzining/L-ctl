package repo

import (
	"github.com/fatih/color"
	"github.com/ningzining/lctl/logic/repo"
	"github.com/spf13/cobra"
)

var dir string
var table string
var style string
var overwrite string

var Cmd = &cobra.Command{
	Use:   "repo",
	Short: "生成仓储层的repo文件",
	Run: func(cmd *cobra.Command, args []string) {
		if err := repo.NewRepo(dir, table, style, overwrite).Generate(); err != nil {
			color.Red("repo文件生成失败, %s\n", err.Error())
			return
		}
	},
}

func init() {
	Cmd.Flags().StringVarP(&dir, "dir", "d", "", "指定文件生成的目录,默认为`./`")
	Cmd.Flags().StringVarP(&table, "table", "t", "", "指定数据库的表名")
	Cmd.Flags().StringVarP(&style, "style", "s", "", "指定生成的文件格式,默认为下划线格式，可改为小驼峰`lCtl`")
	Cmd.Flags().StringVarP(&overwrite, "overwrite", "o", "false", "是否覆盖原有的文件,可改为`true`")

	if err := Cmd.MarkFlagRequired("dir"); err != nil {
		return
	}
	if err := Cmd.MarkFlagRequired("table"); err != nil {
		return
	}
}
