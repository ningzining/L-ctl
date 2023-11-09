package swag

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/ningzining/L-ctl/config"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	urlPrefix = "https://api.apifox.cn/api/v1/projects/"
	urlSuffix = "/import-data"
)

type Swag struct {
	File      string // swagger.json文件路径
	ProjectId string // 项目id
}

func NewSwag(file string, projectId string) *Swag {
	return &Swag{File: file, ProjectId: projectId}
}

func (s *Swag) Upload() error {
	abs, err := filepath.Abs(s.File)
	if err != nil {
		return err
	}
	file, err := os.ReadFile(abs)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s%s%s", urlPrefix, s.ProjectId, urlSuffix)
	ctlConfig, err := config.GetCtlConfig()
	if err != nil {
		return err
	}

	apiFoxReq := ApiFoxReq{
		ImportFormat:        "openapi",
		Data:                string(file),
		ApiOverwriteMode:    "methodAndPath",
		SchemaOverwriteMode: "name",
	}
	apiFoxRes, err := s.post(url, apiFoxReq, ctlConfig.Token)
	if err != nil {
		return err
	}

	if !apiFoxRes.Success {
		return errors.New(fmt.Sprintf("apifox上传失败: %s\n", apiFoxRes.Data))
	}
	color.Green("swagger导入apifox成功\n")
	return nil
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

// Post 根据url获取http的字节流返回
func (s *Swag) post(url string, req ApiFoxReq, token string) (*ApiFoxRes, error) {
	reqData, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}
	reqBody := strings.NewReader(string(reqData))
	client := http.Client{}
	request, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Apifox-Version", "2022-11-16")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	res, err := client.Do(request)
	defer res.Body.Close()
	apiFoxResBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var apiFoxRes ApiFoxRes
	if err := json.Unmarshal(apiFoxResBytes, &apiFoxRes); err != nil {
		return nil, err
	}

	return &apiFoxRes, nil
}
