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
		return errors.New(fmt.Sprintf("mysql url解析异常: %s", err))
	}

	databaseSource := strings.TrimSuffix(arg.Url, "/"+dsn.DBName) + informationSchema
	db, err := sql.NewMysql(databaseSource)
	if err != nil {
		return errors.New(fmt.Sprintf("数据库连接异常: %s", err))
	}

	// 查询目标数据库中所有的表名
	tables, err := model.NewTableRepo(db).GetAllTables(dsn.DBName)
	if err != nil {
		return errors.New(fmt.Sprintf("数据库查询数据异常: %s", err))
	}
	// 获取需要生成model的数据库表表名
	generateTables := getGenerateTables(arg.Tables, tables)
	if len(generateTables) == 0 {
		return errors.New(fmt.Sprintf("不存在需要生成的数据库表"))
	}

	// 获取所有表相关的列信息
	tableMap, err := getTableMap(db, dsn.DBName, generateTables)
	if err != nil {
		return errors.New(fmt.Sprintf("数据库表列信息查询异常: %s", err))
	}

	// 生成目标文件
	err = genModelFromTable(tableMap, arg.Dir)
	if err != nil {
		return errors.New(fmt.Sprintf("生成文件异常: %s", err))
	}
	return nil
}

// 获取需要生成model的数据库表表名
func getGenerateTables(tableArg string, tables []string) []string {
	splitTable := strings.Split(tableArg, ",")
	if len(splitTable) == 0 {
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
func getTableMap(db *gorm.DB, dbName string, tables []string) (map[string]*model.Table, error) {
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

func genModelFromTable(tables map[string]*model.Table, dir string) error {
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
	if dir == "" {
		dir = defaultDir
	}
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	pkg := filepath.Base(abs)
	pkgMap := make(map[string]any)
	pkgMap["pkg"] = pkg

	importsMap := genImport(table)
	typesMap := genTypes(table)

	res := mergeMap(pkgMap, importsMap, typesMap)
	templateDir, err := templateutil.GetModelTemplatePath()
	if err != nil {
		return err
	}
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
	if err = pathutil.Mkdir(dir); err != nil {
		return err
	}

	err = templateutil.SaveTemplateByLocal(templateDir, filePath, res)
	if err != nil {
		return err
	}

	color.Green("文件生成成功: %s", filePath)
	return nil
}

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

func mergeMap(source ...map[string]any) map[string]any {
	res := make(map[string]any)
	for _, t := range source {
		for k, v := range t {
			res[k] = v
		}
	}
	return res
}
