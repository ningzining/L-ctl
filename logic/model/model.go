package model

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	"github.com/ningzining/lctl/util/caseutil"
	"github.com/ningzining/lctl/util/parseutil"
	"github.com/ningzining/lctl/util/pathutil"
	"github.com/ningzining/lctl/util/sqlutil"
	"github.com/ningzining/lctl/util/sqlutil/model"
	"github.com/ningzining/lctl/util/templateutil"
	"gorm.io/gorm"
)

type Model struct {
	Url       string
	Dir       string
	Tables    string
	Overwrite string
	Style     string
}

func NewModel(url, dir, tables, overwrite, style string) *Model {
	return &Model{
		Url:       url,
		Dir:       dir,
		Tables:    tables,
		Overwrite: overwrite,
		Style:     style,
	}
}

// Generate 生成model文件
func (m *Model) Generate() error {
	dsn, err := mysql.ParseDSN(m.Url)
	if err != nil {
		return err
	}

	// 获取mysql系统库的dsn连接
	db, err := sqlutil.NewMysql(m.Url, dsn.DBName)
	if err != nil {
		return err
	}

	// 查询目标数据库中所有的表名
	tables, err := model.NewTableRepo(db).GetAllTables(dsn.DBName)
	if err != nil {
		return err
	}

	// 获取需要生成model的数据库表表名
	tableNames := m.getGenerateTables(m.Tables, tables)
	if len(tableNames) == 0 {
		return err
	}

	// 获取所有表相关的列信息
	tableMap, err := m.getTableStruct(db, dsn.DBName, tableNames)
	if err != nil {
		return err
	}

	// 生成目标文件
	return m.genTemplate(tableMap)
}

// 获取需要生成model的数据库表表名
func (m *Model) getGenerateTables(tableArg string, tables []string) []string {
	splitTable := strings.Split(tableArg, ",")
	// 如果没有输入指定表名，则默认生成全部的表
	if tableArg == "" {
		return tables
	}
	var resTable []string
	for _, t := range splitTable {
		for _, table := range tables {
			if t == table {
				resTable = append(resTable, t)
			}
		}
	}
	return resTable
}

// 获取所有表相关的列信息
func (m *Model) getTableStruct(db *gorm.DB, dbName string, tables []string) (map[string]*model.Table, error) {
	tableMap := make(map[string]*model.Table)
	for _, tableName := range tables {
		table, err := model.NewColumnRepo(db).FindColumns(dbName, tableName)
		if err != nil {
			return nil, err
		}
		tableMap[tableName] = table
	}
	return tableMap, nil
}

// 生成目标model文件
func (m *Model) genTemplate(tables map[string]*model.Table) error {
	for _, item := range tables {
		table, err := parseutil.ConvertTable(item)
		if err != nil {
			return err
		}
		// 生成数据库对应实体的model结构
		if err = m.genModel(table); err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) genModel(table *parseutil.Table) error {
	dirAbs, fileAbs, err := pathutil.CheckAndMkdir(table.TableName, m.Dir, m.Style, m.Overwrite)
	if err != nil {
		color.Red(fmt.Sprintf("当前路径:`%s`已存在目标文件,生成失败", fileAbs))
		return err
	}
	// 获取数据并生成模板文件
	data := getModelTemplateData(dirAbs, table)
	if err := createModelFile(fileAbs, data); err != nil {
		return err
	}

	color.Green("model文件生成成功: %s", fileAbs)
	return nil
}

// 生成model结构体
func getTypes(table *parseutil.Table) map[string]any {
	res := make(map[string]any)
	var fields []map[string]any
	for _, f := range table.Fields {
		field := make(map[string]any)
		field["name"] = caseutil.UpperCamelCase(f.Name)
		field["type"] = f.DataType
		field["tag"] = ""
		field["hasComment"] = f.Comment != ""
		field["comment"] = f.Comment
		fields = append(fields, field)
	}
	res["fields"] = fields
	res["objectName"] = caseutil.UpperCamelCase(table.TableName)
	res["tableName"] = table.TableName
	return res
}

// 生成model模板的数据
func getModelTemplateData(dirAbs string, table *parseutil.Table) map[string]any {
	var importFields []string
	for _, field := range table.Fields {
		importFields = append(importFields, field.DataType)
	}
	return templateutil.MergeMap(templateutil.GetPkg(dirAbs), templateutil.GetImports(importFields), getTypes(table))
}

func createModelFile(filePath string, data map[string]any) (err error) {
	// 获取本地模板文件的路径
	templatePath, err := templateutil.GetModelTemplatePath()
	if err != nil {
		return err
	}

	// 通过本地文件保存模板
	return templateutil.CreateTemplateFile(templatePath, filePath, data)
}
