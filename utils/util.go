package utils

import (
	"os"
	"reflect"
)

func PathExist(path string) bool {
	if _, err := os.Stat(path); err != nil {
		return false
	} else {
		return true
	}
}

func ElemInSlice[T any](elem T, slice []T) bool {
	for i := 0; i < len(slice); i++ {
		if reflect.DeepEqual(elem, slice[i]) {
			return true
		}
	}
	return false
}

func CopyFile(oldFileName string, newFileName string, rewrite bool) error {
	fileData, err := os.ReadFile(oldFileName)
	if err != nil {
		return err
	}
	if rewrite {
		err = os.WriteFile(newFileName, fileData, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		if !PathExist(newFileName) {
			err = os.WriteFile(newFileName, fileData, os.ModePerm)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
