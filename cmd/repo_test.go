package cmd

import (
	"fmt"
	"github.com/ningzining/L-ctl/logic/util/caseutil"
	"net/url"
	"testing"
)

func TestCreateFile(t *testing.T) {
	tableName := "user"
	fileName := fmt.Sprintf("%s.go", caseutil.ToCamelCase(tableName, false))
	filePath, err := url.JoinPath("../tmp", fileName)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	m := make(map[string]interface{})
	m["Name"] = caseutil.ToCamelCase(tableName, true)
	m["TableName"] = tableName
	err = createFile(filePath, m)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	return
}
