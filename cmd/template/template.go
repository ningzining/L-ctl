package template

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "template",
	Short: "操作模板库",
	Long:  "操作模板库，提供了初始化和更新的操作",
}
