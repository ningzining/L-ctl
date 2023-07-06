/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

var dir string
var table string

// repoCmd represents the repo command
var repoCmd = &cobra.Command{
	Use:   "repo",
	Short: "生成仓储层的repo文件",
	Run: func(cmd *cobra.Command, args []string) {
		generate(dir, table)
	},
}

func init() {
	rootCmd.AddCommand(repoCmd)
	repoCmd.Flags().StringVarP(&dir, "dir", "d", "", "指定文件生成的目录")
	repoCmd.Flags().StringVarP(&table, "table", "t", "", "指定数据库的表名(文件名)")
}

func generate(dir string, tableName string) {
	fileName := fmt.Sprintf("%s.go", toCamelCase(tableName, false))
	result, err := url.JoinPath(dir, fileName)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	_, err = os.Stat(result)
	m := make(map[string]interface{})
	m["Name"] = toCamelCase(tableName, true)
	m["TableName"] = tableName
	if err != nil {
		if os.IsNotExist(err) {
			createFile(dir, fileName, m, "template/repo.tpl")
			return
		}
		fmt.Printf("%s\n", err)
	}
	fmt.Printf("目标文件已存在，创建失败: %s\n", result)
	return
}

const (
	UnderLine = '_'
)

func createFile(path string, fileName string, data interface{}, templateFile string) {
	filePath, err := url.JoinPath(path, fileName)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	err = os.MkdirAll(path, 0777)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer file.Close()

	files, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	err = files.Execute(file, data)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	return
}

func toCamelCase(source string, isTitleCase bool) string {
	if len(source) == 0 {
		return source
	}
	var sb strings.Builder
	upper := isTitleCase
	for _, s := range source {
		if s == UnderLine {
			upper = true
		} else if upper {
			toUpper := strings.ToUpper(string(s))
			sb.WriteString(toUpper)
			upper = false
		} else {
			sb.WriteRune(s)
		}
	}
	return sb.String()
}
