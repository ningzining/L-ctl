package pathutil

import (
	"fmt"
	"github.com/ningzining/L-ctl/util/caseutil"
	"os"
	"path/filepath"
)

const (
	CamelCase     = "lCtl"  // 驼峰命名
	UnderLineCase = "l_ctl" // 下划线命名
)

// Exist 判断该路径是否存在
func Exist(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// MkdirIfNotExist 创建文件夹，如果不存在则创建
func MkdirIfNotExist(dir string) error {
	if len(dir) == 0 {
		return nil
	}

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, os.ModePerm)
	}
	return nil
}

// GenFilePath 获取生成文件的路径
func GenFilePath(dirPath, fileName, style string) (string, error) {
	switch style {
	case UnderLineCase:
		fileName = fmt.Sprintf("%s.go", caseutil.ToUnderLineCase(fileName))
	case CamelCase:
		fileName = fmt.Sprintf("%s.go", caseutil.ToCamelCase(fileName, false))
	default:
		fileName = fmt.Sprintf("%s.go", caseutil.ToUnderLineCase(fileName))
	}
	dirAbs, err := filepath.Abs(dirPath)
	filePath := filepath.Join(dirAbs, fileName)
	if err != nil {
		return "", err
	}
	return filePath, nil
}
