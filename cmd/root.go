package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/logic/version"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "L-ctl",
	Short: "L的脚手架工具",
	Long:  `L的脚手架工具，它能够帮助我们自动创建出指定模板的go文件，便于开发`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		color.Red("%s\n", err.Error())
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = fmt.Sprintf("%s %s/%s", version.BuildVersion, runtime.GOOS, runtime.GOARCH)
}
