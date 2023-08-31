package logic

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/util/caseutil"
	"github.com/ningzining/L-ctl/util/pathutil"
	"github.com/ningzining/L-ctl/util/templateutil"
	"path/filepath"
)

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

type RepoGenerateArg struct {
	Dir       string
	Table     string
	Overwrite string
	Style     string
}

// Generate 生成repo文件
func (r *Repo) Generate(arg RepoGenerateArg) error {
	dirAbs, fileAbs, err := pathutil.CheckAndMkdir(arg.Table, arg.Dir, arg.Style, arg.Overwrite)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	// 获取数据并生成模板文件
	data := genRepoTemplateData(dirAbs, arg.Table)
	if err = templateutil.Create(fileAbs, data, templateutil.LocalRepoUrl); err != nil {
		return errors.New(fmt.Sprintf("模板文件生成失败: %s\n", err))
	}

	color.Green("repo文件生成成功: %s", fileAbs)
	return nil
}

func genRepoTemplateData(dirAbs, tableName string) map[string]any {
	pkgMap := map[string]any{
		"pkg": filepath.Base(dirAbs),
	}
	dataMap := map[string]any{
		"name":      caseutil.UpperCamelCase(tableName),
		"tableName": tableName,
	}

	data := templateutil.MergeMap(pkgMap, dataMap)
	return data
}
