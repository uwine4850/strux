package utils

import "fmt"

type ErrPathNotExist struct {
	Path string
}

func (e *ErrPathNotExist) Error() string {
	return fmt.Sprintf("Path %s not exist.", e.Path)
}

type ErrPathAlreadyExist struct {
	Path string
}

func (e *ErrPathAlreadyExist) Error() string {
	return fmt.Sprintf("Path %s already exist.", e.Path)
}
