package commands

import "fmt"

type ErrDbStruxPkgPathNotExist struct {
	RunCommand string
}

func (e *ErrDbStruxPkgPathNotExist) Error() string {
	return fmt.Sprintf("Strux pkg path not exist id database. Please, run command %s.", e.RunCommand)
}

type ErrPackageAlreadyExist struct {
	PkgName string
}

func (e *ErrPackageAlreadyExist) Error() string {
	return fmt.Sprintf("Package %s already exist.", e.PkgName)
}

type ErrCurrentStruxPkgPathNotExist struct {
	PkgPath string
}

func (e *ErrCurrentStruxPkgPathNotExist) Error() string {
	return fmt.Sprintf("Struct pkg path '%s' not exist.", e.PkgPath)
}

type ErrPkgOrConfigNotFound struct {
	PkgName string
}

func (e *ErrPkgOrConfigNotFound) Error() string {
	return fmt.Sprintf("Package or project.toml for '%s' not found.", e.PkgName)
}

type ErrDbStruxPkgPathAlreadyExist struct {
	PkgPath string
}

func (e *ErrDbStruxPkgPathAlreadyExist) Error() string {
	return fmt.Sprintf("The database already contains the path %s.\n"+
		"To change an entry in the database, run the command -i enter/path/ -n to create a new"+
		"directory at the selected path and update the value in the database."+
		"Or run -i enter/path/to/strux_pkg(!) -n to upgrade only the database.", e.PkgPath)
}

type ErrPathMustLeadToStruxPkg struct {
	Path string
}

func (e *ErrPathMustLeadToStruxPkg) Error() string {
	return fmt.Sprintf("The path %s must lead to the strux_pkg directory.", e.Path)
}
