package logic

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
	"path"
	"path/filepath"
	"strings"
)

const (
	decimalImport = "decimal.Decimal"
	timeImport    = "*time.Time"
)

type Model struct{}

func NewModel() *Model {
	return &Model{}
}

const (
	defaultDir        = "./"                  // 默认生成的文件夹目录
	informationSchema = "/information_schema" // 数据库元数据库名
)

type ModelGenerateArg struct {
	Url    string
	Dir    string
	Tables string
}

// Generate 生成model文件
func (m *Model) Generate(arg ModelGenerateArg) error {
	dsn, err := mysql.ParseDSN(arg.Url)
	if err != nil {
		return errors.New(fmt.Sprintf("mysql url解析异常: %s\n", err))
	}

	// 获取mysql系统库的dsn连接
	databaseSource := strings.TrimSuffix(arg.Url, "/"+dsn.DBName) + informationSchema
	db, err := sql.NewMysql(databaseSource)
	if err != nil {
		return errors.New(fmt.Sprintf("数据库连接异常: %s\n", err))
	}

	// 查询目标数据库中所有的表名
	tables, err := model.NewTableRepo(db).GetAllTables(dsn.DBName)
	if err != nil {
		return errors.New(fmt.Sprintf("数据库查询数据异常: %s\n", err))
	}

	// 获取需要生成model的数据库表表名
	tableNames := getGenerateTables(arg.Tables, tables)
	if len(tableNames) == 0 {
		return errors.New(fmt.Sprintf("不存在需要生成的数据库表: %s\n", arg.Tables))
	}

	// 获取所有表相关的列信息
	tableMap, err := getTableStruct(db, dsn.DBName, tableNames)
	if err != nil {
		return errors.New(fmt.Sprintf("数据库表列信息查询异常: %s\n", err))
	}

	// 生成目标文件
	err = genModelTemplate(tableMap, arg.Dir)
	if err != nil {
		return errors.New(fmt.Sprintf("生成文件异常: %s\n", err))
	}
	return nil
}

// 获取需要生成model的数据库表表名
func getGenerateTables(tableArg string, tables []string) []string {
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
func getTableStruct(db *gorm.DB, dbName string, tables []string) (map[string]*model.Table, error) {
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
func genModelTemplate(tables map[string]*model.Table, dir string) error {
	for _, item := range tables {
		table, err := parseutil.ConvertTable(item)
		if err != nil {
			return err
		}
		err = genModel(table, dir)
		if err != nil {
			return err
		}
	}
	return nil
}

func genModel(table *parseutil.Table, dir string) error {
	filePath := path.Join(dir, fmt.Sprintf("%s%s", caseutil.ToUnderLineCase(table.TableName), ".go"))
	// 判断目标文件是否存在
	exist, err := pathutil.Exist(filePath)
	if err != nil {
		return errors.New(fmt.Sprintf("生成模板失败,%s\n", err.Error()))
	}
	if exist {
		return errors.New(fmt.Sprintf("当前路径:`%s`已存在目标文件,生成失败,请选择另外的路径\n", filePath))
	}

	// 创建文件夹
	if dir == "" {
		dir = defaultDir
	}
	dirAbs, err := filepath.Abs(dir)
	if err != nil {
		return errors.New(fmt.Sprintf("获取绝对路径失败: %s\n", err))
	}
	if err = pathutil.MkdirIfNotExist(dirAbs); err != nil {
		return errors.New(fmt.Sprintf("创建目录失败: %s", err))
	}

	// 获取数据并生成模板文件
	data := genModelTemplateData(filePath, table)
	if err = createModelTemplate(filePath, data); err != nil {
		return errors.New(fmt.Sprintf("模板文件生成失败: %s\n", err))
	}

	color.Green("文件生成成功: %s", filePath)
	return nil
}

// 生成model模板的数据
func genModelTemplateData(filePath string, table *parseutil.Table) map[string]any {
	pkgMap := map[string]any{
		"pkg": filepath.Base(filePath),
	}
	// 获取需要import的集合
	importsMap := genImport(table)
	// 获取结构体的集合
	typesMap := genTypes(table)

	data := mergeMap(pkgMap, importsMap, typesMap)
	return data
}

// 生成import结构
func genImport(table *parseutil.Table) map[string]any {
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
func genTypes(table *parseutil.Table) map[string]any {
	res := make(map[string]any)
	var fields []map[string]any
	for _, f := range table.Fields {
		field := make(map[string]any)
		field["name"] = caseutil.ToCamelCase(f.Name, true)
		field["type"] = f.DataType
		field["tag"] = fmt.Sprintf("`gorm:\"column:%s\"`", f.OriginalName)
		field["hasComment"] = f.Comment != ""
		field["comment"] = f.Comment
		fields = append(fields, field)
	}
	res["fields"] = fields
	res["objectName"] = caseutil.ToCamelCase(table.TableName, true)
	return res
}

// 合并所有的map集合
func mergeMap(source ...map[string]any) map[string]any {
	res := make(map[string]any)
	for _, t := range source {
		for k, v := range t {
			res[k] = v
		}
	}
	return res
}

// 创建model模板文件
func createModelTemplate(filePath string, data map[string]any) error {
	// 获取本地模板文件的路径
	templatePath, err := templateutil.GetModelTemplatePath()
	if err != nil {
		return err
	}
	// 通过本地文件保存模板
	if err = templateutil.SaveTemplateByLocal(templatePath, filePath, data); err != nil {
		return err
	}
	return nil
}
