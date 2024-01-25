package model

import (
	"fmt"
	"testing"
)

func TestModel_Generate(t *testing.T) {
	err := NewModel("root:root@tcp(127.0.0.1:3306)/test", "./cache", "", "true", "").Generate()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
}
