package model

import (
	"gorm.io/gorm"
)

type ColumnRepo struct {
	db *gorm.DB
}

type (
	Table struct {
		DbName    string    // 数据库名
		TableName string    // 表名
		Columns   []*Column // 列名
	}
	Column struct {
		*DbColumn
		Index []*DbIndex
	}
	DbColumn struct {
		Name            string `gorm:"column:COLUMN_NAME"`
		DataType        string `gorm:"column:DATA_TYPE"`
		ColumnType      string `gorm:"column:COLUMN_TYPE"`
		Extra           string `gorm:"column:EXTRA"`
		Comment         string `gorm:"column:COLUMN_COMMENT"`
		ColumnDefault   string `gorm:"column:COLUMN_DEFAULT"`
		IsNullAble      string `gorm:"column:IS_NULLABLE"`
		OrdinalPosition int    `gorm:"column:ORDINAL_POSITION"`
	}

	// DbIndex defines index of columns in information_schema.statistic
	DbIndex struct {
		IndexName  string `gorm:"column:INDEX_NAME"`
		NonUnique  int    `gorm:"column:NON_UNIQUE"`
		SeqInIndex int    `gorm:"column:SEQ_IN_INDEX"`
	}
)

type IColumnRepo interface {
	FindColumns(dbName, tableName string) (*Table, error)
}

func NewColumnRepo(db *gorm.DB) IColumnRepo {
	return &ColumnRepo{db: db}
}

func (i *ColumnRepo) TableName() string {
	return "COLUMNS"
}

func (i *ColumnRepo) FindColumns(dbName, tableName string) (*Table, error) {
	var columns []*DbColumn
	err := i.db.Table(i.TableName()).
		Select("COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,EXTRA,COLUMN_COMMENT,COLUMN_DEFAULT,IS_NULLABLE,ORDINAL_POSITION").
		Where("TABLE_SCHEMA = ?", dbName).
		Where("TABLE_NAME = ?", tableName).
		Order("ORDINAL_POSITION asc").
		Find(&columns).
		Error
	if err != nil {
		return nil, err
	}
	var list []*Column
	for _, item := range columns {
		index, err := NewStatisticRepo(i.db).FindIndex(dbName, tableName, item.Name)
		if err != nil {
			return nil, err
		}
		if len(index) > 0 {
			var indexList []*DbIndex
			for _, temp := range index {
				indexList = append(indexList, temp)
			}
			list = append(list, &Column{
				DbColumn: item,
				Index:    indexList,
			})
		} else {
			list = append(list, &Column{
				DbColumn: item,
			})
		}
	}
	table := &Table{
		DbName:    dbName,
		TableName: tableName,
		Columns:   list,
	}
	return table, nil
}
