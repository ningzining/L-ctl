package logic

import (
	"fmt"
	"github.com/ningzining/L-ctl/cache"
	"github.com/ningzining/L-ctl/sql"
	"testing"
)

func TestModel_Generate(t *testing.T) {
	arg := ModelGenerateArg{
		Url:    "root:root@tcp(127.0.0.1:3306)/test",
		Dir:    "./cache",
		Tables: "",
	}
	err := NewModel().Generate(arg)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}

func TestModel_Auto(t *testing.T) {
	mysql, err := sql.NewMysql("root:root@tcp(127.0.0.1:3306)/test")
	if err != nil {
		return
	}
	mysql.AutoMigrate(&cache.Sysusers{})
}
