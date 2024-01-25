package apifox

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	urlPrefix = "https://api.apifox.cn/api/v1/projects/"
	urlSuffix = "/import-data"
)

type ImportDataReq struct {
	ImportFormat        string `json:"importFormat"`
	Data                string `json:"data"`
	ApiOverwriteMode    string `json:"apiOverwriteMode"`
	SchemaOverwriteMode string `json:"schemaOverwriteMode"`
}

type ImportDataRes struct {
	Data    any  `json:"data"`
	Success bool `json:"success"`
}

type Client struct {
	Token string
}

func NewClient(token string) *Client {
	return &Client{
		Token: token,
	}
}

// ImportData apifox导入数据
func (c *Client) ImportData(projectId string, req ImportDataReq) (*ImportDataRes, error) {
	reqData, err := json.Marshal(&req)
	if err != nil {
		return nil, err
	}

	resBytes, err := c.post(c.getImportDataUrl(projectId), reqData)
	if err != nil {
		return nil, err
	}

	var apiFoxRes ImportDataRes
	if err := json.Unmarshal(resBytes, &apiFoxRes); err != nil {
		return nil, err
	}

	return &apiFoxRes, nil
}

func (c *Client) getImportDataUrl(projectId string) string {
	return fmt.Sprintf("%s%s%s", urlPrefix, projectId, urlSuffix)
}

// post apifox获取http的的返回
func (c *Client) post(url string, req []byte) ([]byte, error) {
	reqBody := strings.NewReader(string(req))

	request, err := http.NewRequest("POST", url, reqBody)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Apifox-Version", "2022-11-16")
	request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	client := http.Client{}
	res, err := client.Do(request)
	defer res.Body.Close()

	apiFoxResBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return apiFoxResBytes, nil
}
