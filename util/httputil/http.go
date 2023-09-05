package httputil

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Get 根据url获取http的字节流返回
func Get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// Post 根据url获取http的字节流返回
func Post(url string, req interface{}) ([]byte, error) {
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
	request.Header.Set("Authorization", "Bearer APS-a0Vv1MWiWXQ4UUDs1OnwmaojLGXVIR2Z")
	res, err := client.Do(request)
	defer res.Body.Close()
	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
