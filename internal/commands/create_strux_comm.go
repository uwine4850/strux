package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strux/internal/db"
	"strux/internal/filegen"
	"strux/utils"
)

type CommandInterface interface {
	OnFinish()
}

type CreateCommand struct {
	Create string `short:"crt" long:"create" block:"1"`
}

// ExecCreate creating a new package.
// The package is created in the place where the application was previously initialized.
// In addition to the package directory, a base settings file is also created.
func (cc *CreateCommand) ExecCreate(name string) []string {
	struxPkgPath, err := db.GetStruxPkgPathValue()
	if err != nil {
		fmt.Println("Please, run command --init path/to/strux/pkg")
		panic(err)
	}
	if struxPkgPath != "" && utils.PathExist(struxPkgPath) {
		if utils.PathExist(filepath.Join(struxPkgPath, name)) {
			fmt.Println(fmt.Sprintf("Package %s already exist.", name))
			return []string{}
		}
		if err := os.Mkdir(filepath.Join(struxPkgPath, name), os.ModePerm); err != nil {
			panic(err)
		} else {
			gen := filegen.TomlGenerator{
				PkgName:     name,
				Version:     "0.0.1",
				Description: "",
				PkgPath:     filepath.Join(struxPkgPath, name),
			}
			gen.Generate()
			fmt.Println(fmt.Sprintf("Package %s created successfully.", name))
		}
	} else {
		fmt.Println(fmt.Sprintf("Struct pkg path '%s' not exist.", struxPkgPath))
		fmt.Println("Create directories manually or run command")
	}
	return []string{}
}

func (cc *CreateCommand) OnFinish() {

}
