package parseutil

import (
	"errors"
	"fmt"
	"github.com/ningzining/L-ctl/util/caseutil"
	"github.com/ningzining/L-ctl/util/sqlutil/model"
	"sort"
	"strings"
)

var commonMysqlDataTypeMapString = map[string]string{
	// For consistency, all integer types are converted to int64
	// bool
	"bool":    "bool",
	"_bool":   "pq.BoolArray",
	"boolean": "bool",
	// number
	"tinyint":   "bool",
	"smallint":  "int64",
	"mediumint": "int64",
	"int":       "int64",
	"int1":      "int64",
	"int2":      "int64",
	"_int2":     "pq.Int64Array",
	"int3":      "int64",
	"int4":      "int64",
	"_int4":     "pq.Int64Array",
	"int8":      "int64",
	"_int8":     "pq.Int64Array",
	"integer":   "int64",
	"_integer":  "pq.Int64Array",
	"bigint":    "int64",
	"float":     "float64",
	"float4":    "float64",
	"_float4":   "pq.Float64Array",
	"float8":    "float64",
	"_float8":   "pq.Float64Array",
	"double":    "float64",
	"decimal":   "decimal.Decimal",
	"dec":       "float64",
	"fixed":     "float64",
	"real":      "float64",
	"bit":       "byte",
	// date & time
	"date":      "*time.Time",
	"datetime":  "*time.Time",
	"timestamp": "*time.Time",
	"time":      "string",
	"year":      "int64",
	// string
	"linestring":      "string",
	"multilinestring": "string",
	"nvarchar":        "string",
	"nchar":           "string",
	"char":            "string",
	"bpchar":          "string",
	"_char":           "pq.StringArray",
	"character":       "string",
	"varchar":         "string",
	"_varchar":        "pq.StringArray",
	"binary":          "string",
	"bytea":           "string",
	"longvarbinary":   "string",
	"varbinary":       "string",
	"tinytext":        "string",
	"text":            "string",
	"_text":           "pq.StringArray",
	"mediumtext":      "string",
	"longtext":        "string",
	"enum":            "string",
	"set":             "string",
	"json":            "string",
	"jsonb":           "string",
	"blob":            "string",
	"longblob":        "string",
	"mediumblob":      "string",
	"tinyblob":        "string",
	"ltree":           "[]byte",
}

const (
	primary = "PRIMARY"
)

type Table struct {
	TableName  string
	DbName     string
	Fields     []*Field
	PrimaryKey string
}

type Field struct {
	OriginalName    string // 数据库列名，用于生成gorm的tag
	Name            string // 字段名
	DataType        string // 数据类型
	Comment         string // 注释
	OrdinalPosition int    // 排序顺序
}

// ConvertTable 转换数据库相关的字段为go结构体相关字段
func ConvertTable(table *model.Table) (*Table, error) {
	var resTable Table
	resTable.TableName = table.TableName
	resTable.DbName = table.DbName

	var fields []*Field
	for _, column := range table.Columns {
		for _, index := range column.Index {
			if index.IndexName == primary {
				resTable.PrimaryKey = column.Name
			}
		}
		goType, err := ConvertStringDataType(column.DataType)
		if err != nil {
			return nil, err
		}
		field := &Field{
			OriginalName:    column.Name,
			Name:            caseutil.UpperCamelCase(column.Name),
			DataType:        goType,
			Comment:         strings.NewReplacer("\r", "", "\n", "").Replace(column.Comment),
			OrdinalPosition: column.OrdinalPosition,
		}
		fields = append(fields, field)
	}
	sort.Slice(fields, func(i, j int) bool {
		return fields[i].OrdinalPosition < fields[j].OrdinalPosition
	})
	resTable.Fields = fields
	return &resTable, nil
}

// ConvertStringDataType 转换mysql的数据类型为go的数据类型
func ConvertStringDataType(sourceType string) (goType string, err error) {
	tp, ok := commonMysqlDataTypeMapString[strings.ToLower(sourceType)]
	if !ok {
		return "", errors.New(fmt.Sprintf("不支持目标mysql列类型: %s", sourceType))
	}
	return tp, nil
}
