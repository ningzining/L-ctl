package pathutil

import (
	"errors"
	"fmt"
	"github.com/ningzining/L-ctl/util/caseutil"
	"os"
	"path/filepath"
)

const (
	LowerCamelCase = "lCtl"  // 驼峰命名
	UnderLineCase  = "l_ctl" // 下划线命名
)

const (
	defaultDir = "./"
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
		fileName = fmt.Sprintf("%s.go", caseutil.UnderLineCase(fileName))
	case LowerCamelCase:
		fileName = fmt.Sprintf("%s.go", caseutil.LowerCamelCase(fileName))
	default:
		fileName = fmt.Sprintf("%s.go", caseutil.UnderLineCase(fileName))
	}
	dirAbs, err := filepath.Abs(dirPath)
	filePath := filepath.Join(dirAbs, fileName)
	if err != nil {
		return "", err
	}
	return filePath, nil
}

// CheckAndMkdir 检查文件是否存在并且创建文件夹
func CheckAndMkdir(tableName, argDir, fileStyle, overwrite string) (dirAbs, fileAbs string, err error) {
	// 创建文件夹
	dir := argDir
	if dir == "" {
		dir = defaultDir
	}
	dirAbs, err = filepath.Abs(dir)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("获取绝对路径失败: %s\n", err))
	}
	if err = MkdirIfNotExist(dirAbs); err != nil {
		return "", "", errors.New(fmt.Sprintf("创建目录失败: %s", err))
	}

	// 判断目标文件是否存在
	filePath, err := GenFilePath(dirAbs, tableName, fileStyle)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("当前路径:`%s`已存在目标文件,生成失败", fileAbs))
	}
	exist, err := Exist(filePath)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("生成模板失败,%s\n", err.Error()))
	}
	if exist && overwrite != "true" {
		return "", "", errors.New(fmt.Sprintf("当前路径:`%s`已存在目标文件,生成失败", filePath))
	}
	return dirAbs, filePath, nil
}
