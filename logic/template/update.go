package template

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/util/gitutil"
	"github.com/ningzining/L-ctl/util/templateutil"
)

// Update 更新模板文件
func (t *Template) Update() error {
	templateDir, err := templateutil.GetTemplateDir()
	if err != nil {
		return err
	}

	output, err := gitutil.Pull(templateDir)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	color.Green("模板库更新成功: %s", templateDir)
	return nil
}
