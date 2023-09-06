package model

import "gorm.io/gorm"

type StatisticRepo struct {
	db *gorm.DB
}

type IStatisticRepo interface {
	FindIndex(database, table, column string) ([]*DbIndex, error)
}

func NewStatisticRepo(db *gorm.DB) IStatisticRepo {
	return &StatisticRepo{db: db}
}

func (s *StatisticRepo) TableName() string {
	return "STATISTICS"
}

func (s *StatisticRepo) FindIndex(database, table, column string) ([]*DbIndex, error) {
	var reply []*DbIndex
	err := s.db.Table(s.TableName()).
		Select("INDEX_NAME,NON_UNIQUE,SEQ_IN_INDEX").
		Where("TABLE_SCHEMA = ?", database).
		Where("TABLE_NAME = ?", table).
		Where("COLUMN_NAME = ?", column).
		Find(&reply).Error
	if err != nil {
		return nil, err
	}
	return reply, nil
}
