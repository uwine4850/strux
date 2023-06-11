package commands

import (
	"fmt"
	"reflect"
	"strux/internal/db"
)

var fieldInfo = map[string]string{
	"Pkg": "Getting strux_pkg path.",
	"New": "New info",
}

type InfoCommand struct {
	Info   string `short:"i" long:"info" block:"1"`
	Pkg    string `short:"p" long:"pkg"`
	New    string `short:"n" long:"new"`
	isInfo bool
	isPkg  bool
	isNew  bool
}

func (inf *InfoCommand) ExecInfo() []string {
	inf.isInfo = true
	inf.isPkg = false
	inf.isNew = false
	return []string{inf.Pkg}
}

func (inf *InfoCommand) ExecPkg() []string {
	inf.isPkg = true
	value, err := db.GetStruxPkgPathValue()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	return []string{}
}

func (inf *InfoCommand) ExecNew() []string {
	inf.isNew = true
	fmt.Println("NEW")
	return []string{}
}

func (inf *InfoCommand) OnFinish() {
	if inf.isInfo {
		if !inf.isPkg && !inf.isNew {
			inf.getInfo()
		}
	}
}

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
