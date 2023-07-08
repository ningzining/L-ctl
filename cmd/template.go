package cmd

import (
	"github.com/spf13/cobra"
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "操作模板库",
	Long:  "操作模板库，提供了初始化和更新的操作",
}

func init() {
	rootCmd.AddCommand(templateCmd)
}
