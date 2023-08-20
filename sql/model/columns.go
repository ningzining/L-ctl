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
		Name            string `gorm:"column:COLUMN_NAME"`
		DataType        string `gorm:"column:DATA_TYPE"`
		ColumnType      string `gorm:"column:COLUMN_TYPE"`
		Extra           string `gorm:"column:EXTRA"`
		Comment         string `gorm:"column:COLUMN_COMMENT"`
		ColumnDefault   string `gorm:"column:COLUMN_DEFAULT"`
		IsNullAble      string `gorm:"column:IS_NULLABLE"`
		OrdinalPosition int    `gorm:"column:ORDINAL_POSITION"`
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
	var column []*Column
	err := i.db.Table(i.TableName()).
		Select("COLUMN_NAME,DATA_TYPE,COLUMN_TYPE,EXTRA,COLUMN_COMMENT,COLUMN_DEFAULT,IS_NULLABLE,ORDINAL_POSITION").
		Where("TABLE_SCHEMA = ?", dbName).
		Where("TABLE_NAME = ?", tableName).
		Order("ORDINAL_POSITION asc").
		Find(&column).
		Error
	if err != nil {
		return nil, err
	}
	table := &Table{
		DbName:    dbName,
		TableName: tableName,
		Columns:   column,
	}
	return table, nil
}
