package repo

import (
	"gorm.io/gorm"
)

type {{.Name}}Repo struct {
	db *gorm.DB
}

type I{{.Name}}Repo interface {
	Create(po *model.{{.Name}}) error
	CreateInBatches(list []*model.{{.Name}}) error
	Delete(id int64) error
	Update(po *model.{{.Name}}) error
	FindById(id int64) (*model.{{.Name}}, error)
	FindAll() ([]*model.{{.Name}}, error)
}

func New{{.Name}}(db *gorm.DB) I{{.Name}}Repo {
	return &{{.Name}}Repo{db: db}
}

func (r *{{.Name}}Repo) TableName() string {
	return "{{.TableName}}"
}

func (r *{{.Name}}Repo) Create(po *model.{{.Name}}) error {
	return r.db.Table(r.TableName()).Create(&po).Error
}

func (r *{{.Name}}Repo) CreateInBatches(list []*model.{{.Name}}) error {
	return r.db.Table(r.TableName()).CreateInBatches(list, len(list)).Error
}

func (r *{{.Name}}Repo) Delete(id int64) error {
	return r.db.Table(r.TableName()).Where("id = ?", id).Update("deleted", 1).Error
}

func (r *{{.Name}}Repo) Update(po *model.{{.Name}}) error {
	return r.db.Table(r.TableName()).Where("id = ?", po.Id).Updates(&po).Error
}

func (r *{{.Name}}Repo) FindById(id int64) (po *model.{{.Name}}, err error) {
	err = r.db.Table(r.TableName()).Where("id = ?", id).Where("deleted = 0").First(&po).Error
	return
}

func (r *{{.Name}}Repo) FindAll() (list []*model.{{.Name}}, err error) {
	err = r.db.Table(r.TableName()).Where("deleted = 0").Order("id asc").Find(&list).Error
	return
}