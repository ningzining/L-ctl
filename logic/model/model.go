package model

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-sql-driver/mysql"
	"github.com/ningzining/L-ctl/sql"
	"github.com/ningzining/L-ctl/sql/model"
	"github.com/ningzining/L-ctl/util/caseutil"
	"github.com/ningzining/L-ctl/util/parseutil"
	"github.com/ningzining/L-ctl/util/pathutil"
	"github.com/ningzining/L-ctl/util/templateutil"
	"gorm.io/gorm"
	"path/filepath"
	"strings"
)

const (
	decimalImport = "decimal.Decimal"
	timeImport    = "time.Time"
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
	db, err := sql.NewMysql(m.Url, dsn.DBName)
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
	err = m.genTemplate(tableMap)
	if err != nil {
		return err
	}
	return nil
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
		err = m.genModel(table)
		if err != nil {
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
	data := m.genModelTemplateData(dirAbs, table)
	if err = templateutil.Create(fileAbs, data, templateutil.LocalModelUrl); err != nil {
		return errors.New(fmt.Sprintf("模板文件生成失败: %s\n", err))
	}

	color.Green("model文件生成成功: %s", fileAbs)
	return nil
}

// 生成model模板的数据
func (m *Model) genModelTemplateData(dirAbs string, table *parseutil.Table) map[string]any {
	pkgMap := map[string]any{
		"pkg": filepath.Base(dirAbs),
	}
	// 获取需要import的集合
	importsMap := m.genImport(table)
	// 获取结构体的集合
	typesMap := m.genTypes(table)

	data := templateutil.MergeMap(pkgMap, importsMap, typesMap)
	return data
}

// 生成import结构
func (m *Model) genImport(table *parseutil.Table) map[string]any {
	importsMap := make(map[string]any)
	for _, column := range table.Fields {
		if column.DataType == timeImport {
			importsMap["time"] = true
		}
		if column.DataType == decimalImport {
			importsMap["decimal"] = true
		}
	}
	return importsMap
}

// 生成model结构体
func (m *Model) genTypes(table *parseutil.Table) map[string]any {
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
