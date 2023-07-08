package cmd

import (
	"errors"
	"fmt"
	"github.com/ningzining/L-ctl/logic/util/caseutil"
	"github.com/ningzining/L-ctl/logic/util/httputil"
	"github.com/ningzining/L-ctl/logic/util/pathutil"
	"net/url"
	"os"
	"text/template"

	"github.com/spf13/cobra"
)

var dir string
var table string

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "生成仓储层的repo文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		return generate(dir, table)
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.Flags().StringVarP(&dir, "dir", "d", "", "指定文件生成的目录")
	repoCmd.Flags().StringVarP(&table, "table", "t", "", "指定数据库的表名(文件名)")
	err := repoCmd.MarkFlagRequired("dir")
	if err != nil {
		return
	}
	err = repoCmd.MarkFlagRequired("table")
	if err != nil {
		return
	}
}

func generate(dirPath string, tableName string) error {
	fileName := fmt.Sprintf("%s.go", caseutil.ToCamelCase(tableName, false))
	filePath, err := url.JoinPath(dirPath, fileName)
	if err != nil {
		return err
	}
	// 判断目标文件是否已经存在
	exist, err := pathutil.Exist(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("生成模板失败,%s\n", err.Error()))
	}
	if exist {
		return errors.New(fmt.Sprintf("当前路径:`%s`已存在目标文件,生成失败,请选择另外的路径\n", filePath))
	}
	// 创建文件夹
	err = pathutil.Mkdir(dir)
	if err != nil {
		return err
	}
	// 新建文件
	m := make(map[string]interface{})
	m["Name"] = caseutil.ToCamelCase(tableName, true)
	m["TableName"] = tableName
	err = createFile(filePath, m, "https://raw.githubusercontent.com/ningzining/L-ctl-template/main/repo/repo.tpl")
	if err != nil {
		return err
	}
	return nil
}

func createFile(filePath string, data interface{}, templateUrl string) error {
	// 获取数据
	templateData, err := httputil.Get(templateUrl)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 解析tpl模板文件
	files, err := template.New("temp.tpl").Parse(string(templateData))
	if err != nil {
		return err
	}
	// 渲染数据到模板
	err = files.Execute(file, data)
	if err != nil {
		return err
	}
	return nil
}
