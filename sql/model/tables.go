package model

import "gorm.io/gorm"

type TableRepo struct {
	db *gorm.DB
}

type ITableRepo interface {
	GetAllTables(database string) ([]string, error)
}

func NewTableRepo(db *gorm.DB) ITableRepo {
	return &TableRepo{db: db}
}

func (i *TableRepo) TableName() string {
	return "TABLES"
}

func (i *TableRepo) GetAllTables(database string) ([]string, error) {
	var tables []string
	err := i.db.Table(i.TableName()).Select("TABLE_NAME").Where("TABLE_SCHEMA = ?", database).Find(&tables).Error
	if err != nil {
		return nil, err
	}
	return tables, nil
}
