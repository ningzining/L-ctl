package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/util/httputil"
	"github.com/ningzining/L-ctl/util/templateutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	urlPrefix = "https://api.apifox.cn/api/v1/projects/"
	urlSuffix = "/import-data"
)

type Swag struct{}

func NewSwag() *Swag {
	return &Swag{}
}

type SwagGenerateArg struct {
	File      string
	ProjectId string
}
type ApiFoxReq struct {
	ImportFormat        string `json:"importFormat"`
	Data                string `json:"data"`
	ApiOverwriteMode    string `json:"apiOverwriteMode"`
	SchemaOverwriteMode string `json:"schemaOverwriteMode"`
}
type ApiFoxRes struct {
	Data    interface{} `json:"data"`
	Success bool        `json:"success"`
}

func (s *Swag) Upload(arg SwagGenerateArg) error {
	abs, err := filepath.Abs(arg.File)
	if err != nil {
		return err
	}
	file, err := os.ReadFile(abs)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%s%s", urlPrefix, arg.ProjectId, urlSuffix)

	apiFoxReq := ApiFoxReq{
		ImportFormat:        "openapi",
		Data:                string(file),
		ApiOverwriteMode:    "methodAndPath",
		SchemaOverwriteMode: "name",
	}
	configFile, err := templateutil.GenerateConfigFile()
	if err != nil {
		return err
	}
	tokenBytes, err := os.ReadFile(configFile)
	if err != nil {
		return err
	}

	token := string(tokenBytes)
	token = strings.ReplaceAll(token, "\n", "")
	token = strings.ReplaceAll(token, "\t", "")
	token = strings.ReplaceAll(token, " ", "")
	if token == "" {
		return errors.New("apifox的token不存在，请配置到指定文件中")
	}
	apiFoxResData, err := httputil.Post(url, apiFoxReq, token)
	if err != nil {
		return err
	}
	var apiFoxRes ApiFoxRes
	err = json.Unmarshal(apiFoxResData, &apiFoxRes)
	if err != nil {
		return err
	}
	if apiFoxRes.Success {
		color.Green("apifox重新生成成功")
	}
	return nil
}
