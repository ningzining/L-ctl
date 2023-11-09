package model

import (
	"fmt"
	"github.com/ningzining/L-ctl/cache"
	"github.com/ningzining/L-ctl/util/sqlutil"
	"testing"
)

func TestModel_Generate(t *testing.T) {
	err := NewModel("root:root@tcp(127.0.0.1:3306)/test", "./cache", "", "true", "").Generate()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func TestModel_Auto(t *testing.T) {
	mysql, err := sqlutil.NewMysql("root:root@tcp(127.0.0.1:3306)", "test")
	if err != nil {
		return
	}
	mysql.AutoMigrate(cache.Test{})
}
