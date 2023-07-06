package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "L-ctl",
	Short:   "L的脚手架工具",
	Long:    `L的脚手架工具，它能够帮助我们自动创建出指定模板的go文件，便于开发`,
	Version: version,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

}
