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

type ErrThisIsNotADir struct {
	DirPath string
}

func (e *ErrThisIsNotADir) Error() string {
	return fmt.Sprintf("The path %s does not lead to a directory.", e.DirPath)
}

type ErrDirNotCreated struct {
	DirPath string
}

func (e *ErrDirNotCreated) Error() string {
	return fmt.Sprintf("The %s not created.", e.DirPath)
}
