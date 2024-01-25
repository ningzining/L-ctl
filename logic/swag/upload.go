package swag

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/config"
	"github.com/ningzining/L-ctl/logic/swag/apifox"
)

func (s *Swag) Upload() error {
	abs, err := filepath.Abs(s.File)
	if err != nil {
		return err
	}
	file, err := os.ReadFile(abs)
	if err != nil {
		return err
	}

	ctlConfig, err := config.GetCtlConfig()
	if err != nil {
		return err
	}

	apiFoxReq := apifox.ImportDataReq{
		ImportFormat:        "openapi",
		Data:                string(file),
		ApiOverwriteMode:    "methodAndPath",
		SchemaOverwriteMode: "name",
	}
	apiFoxRes, err := apifox.NewClient(ctlConfig.Token).ImportData(s.ProjectId, apiFoxReq)
	if err != nil {
		return err
	}

	if !apiFoxRes.Success {
		return errors.New(fmt.Sprintf("%v", apiFoxRes.Data))
	}
	color.Green("swagger文件导入apifox成功")
	return nil
}
