package commands

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"reflect"
	"strux/internal/config"
	"strux/internal/db"
	"strux/internal/filegen"
	"strux/utils"
)

var fieldInfo = map[string]string{
	"Pkg": "Getting strux_pkg path.",
	"New": "New info",
}

type InfoCommand struct {
	Info      string `short:"info" long:"info" block:"1"`
	Path      string `short:"p" long:"path"`
	Package   string `short:"pkg" long:"package"`
	isInfo    bool
	isPath    bool
	isPackage bool
}

func (inf *InfoCommand) ExecInfo() []string {
	inf.isInfo = true
	inf.isPath = false
	inf.isPackage = false
	return []string{inf.Path, inf.Package}
}

// ExecPath shows struct_pkg path.
func (inf *InfoCommand) ExecPath() []string {
	inf.isPath = true
	value, err := db.GetStruxPkgPathValue()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	return []string{}
}

// ExecPackage shows information from the project.toml file for the selected package.
func (inf *InfoCommand) ExecPackage(pkgName string) []string {
	inf.isPackage = true
	struxPkgPath, err := db.GetStruxPkgPathValue()
	if err != nil {
		panic(err)
	}
	pkgConfPath := filepath.Join(struxPkgPath, pkgName, config.ProjectConfName)
	if utils.PathExist(pkgConfPath) {
		var fd filegen.FileData
		file, err := os.ReadFile(pkgConfPath)
		if err != nil {
			return nil
		}
		// get/set data
		_, err = toml.Decode(string(file), &fd)
		if err != nil {
			return nil
		}
		inf.printInfoPkgValue(&fd)
	} else {
		fmt.Println(fmt.Sprintf("Package or project.toml for '%s' not found.", pkgName))
	}
	return []string{}
}

func (inf *InfoCommand) OnFinish() {
	if inf.isInfo {
		if !inf.isPath && !inf.isPackage {
			inf.getInfo()
		}
	}
}

// printInfoPkgValue prints in formatted form all package data from the configuration file.
func (inf *InfoCommand) printInfoPkgValue(pkg *filegen.FileData) {
	p := reflect.ValueOf(pkg)
	for i := 0; i < p.Elem().Type().NumField(); i++ {
		fieldName := p.Elem().Type().Field(i).Name
		fieldValue := reflect.Indirect(p).Field(i).String()
		fmt.Println(fmt.Sprintf("%s: %s", fieldName, fieldValue))
	}
}

// getInfo displays information about all current commands.
func (inf *InfoCommand) getInfo() {
	f := reflect.ValueOf(inf).Elem()
	for i := 0; i < f.NumField(); i++ {
		field := f.Type().Field(i)
		if _, ok := field.Tag.Lookup("short"); !ok {
			continue
		}
		if _, ok := field.Tag.Lookup("long"); !ok {
			continue
		}
		for fName, fInfo := range fieldInfo {
			if fName == field.Name {
				fmt.Println(fmt.Sprintf("-%s, --%s | %s",
					field.Tag.Get("short"),
					field.Tag.Get("long"),
					fInfo))
			}
		}
	}
}
