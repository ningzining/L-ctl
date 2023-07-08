package pathutil

import (
	"os"
)

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

func Mkdir(path string) error {
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return err
	}
	return nil
}
