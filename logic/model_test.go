package logic

import (
	"fmt"
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
