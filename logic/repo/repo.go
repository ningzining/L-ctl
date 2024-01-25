package repo

import (
	"github.com/fatih/color"
	"github.com/ningzining/lctl/util/caseutil"
	"github.com/ningzining/lctl/util/pathutil"
	"github.com/ningzining/lctl/util/templateutil"
)

type Repo struct {
	Dir       string
	Table     string
	Overwrite string
	Style     string
}

func NewRepo(dir, table, overwrite, style string) *Repo {
	return &Repo{
		Dir:       dir,
		Table:     table,
		Overwrite: overwrite,
		Style:     style,
	}
}

// Generate 生成repo文件
func (r *Repo) Generate() error {
	dirAbs, fileAbs, err := pathutil.CheckAndMkdir(r.Table, r.Dir, r.Style, r.Overwrite)
	if err != nil {
		color.Red(err.Error())
		return err
	}
	// 获取数据并生成模板文件
	data := getRepoTemplateData(dirAbs, r.Table)
	if err := createRepoFile(fileAbs, data); err != nil {
		return err
	}

	color.Green("repo文件生成成功: %s", fileAbs)
	return nil
}

func getRepoTemplateData(dirAbs, tableName string) map[string]any {
	dataMap := map[string]any{
		"name":      caseutil.UpperCamelCase(tableName),
		"tableName": tableName,
	}

	return templateutil.MergeMap(templateutil.GetPkg(dirAbs), dataMap)
}

func createRepoFile(filePath string, data map[string]any) (err error) {
	// 获取repo模板文件路径
	templatePath, err := templateutil.GetRepoTemplatePath()
	if err != nil {
		return err
	}
	// 创建模板文件
	return templateutil.CreateTemplateFile(templatePath, filePath, data)
}
