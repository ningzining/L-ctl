package model

import (
	"github.com/ningzining/L-ctl/logic"
	"github.com/spf13/cobra"
)

var url string
var dir string
var tables string
var overwrite string
var style string

var Cmd = &cobra.Command{
	Use:   "model",
	Short: "生成数据库层的model文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		arg := logic.ModelGenerateArg{
			Url:       url,
			Dir:       dir,
			Tables:    tables,
			Overwrite: overwrite,
			Style:     style,
		}
		return logic.NewModel().Generate(arg)
	},
}

func init() {
	Cmd.Flags().StringVarP(&dir, "dir", "d", "", "指定文件生成的目录")
	Cmd.Flags().StringVarP(&tables, "table", "t", "", "指定数据库的表名,多个表用`,`隔开")
	Cmd.Flags().StringVarP(&url, "url", "u", "", "指定数据库dsn的连接")
	Cmd.Flags().StringVarP(&style, "style", "s", "", "设置文件名的格式（默认下划线，可改为驼峰）")
	Cmd.Flags().StringVarP(&overwrite, "overwrite", "o", "false", "是否覆盖原有的文件")
	if err := Cmd.MarkFlagRequired("url"); err != nil {
		return
	}
}
