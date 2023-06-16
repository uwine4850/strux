package commands

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"strux/internal/config"
	"strux/internal/db"
	"strux/utils"
)

type InitCommand struct {
	Init       string `short:"i" long:"init" block:"1"`
	New        string `short:"n" long:"new"`
	Database   string `short:"db" long:"database"`
	Path       string
	isNew      bool
	isDatabase bool
}

func (ic *InitCommand) ExecInit(path string) []string {
	ic.Path = path
	if !utils.PathExist(path) {
		err := &utils.ErrPathNotExist{Path: ic.Path}
		panic(err)
	}
	return []string{ic.New, ic.Database}
}

func (ic *InitCommand) ExecNew() []string {
	ic.isNew = true
	return []string{}
}

func (ic *InitCommand) ExecDatabase() []string {
	ic.isDatabase = true
	return []string{}
}

func (ic *InitCommand) OnFinish() {
	if ic.isNew {
		ic.acceptedNewPath()
	} else if ic.isDatabase {
		ic.updateStructPkgDb()
	} else {
		ic.acceptedPath()
	}
}

// acceptedNewPath create a strux_pkg directory at the specified path.
// Updates the corresponding entry in the database.
func (ic *InitCommand) acceptedNewPath() {
	ic.createStruxPkgFolder(func() {
		err := db.CreateDbTable()
		if err != nil {
			panic(err.Error())
		} else {
			dirpath := filepath.Join(ic.Path, config.StruxPkgName)
			_, err := db.ExecStruxDbQuery(fmt.Sprintf("UPDATE strux SET strux_pkg_path = '%s' WHERE id = 1", dirpath))
			if err != nil {
				panic(err)
			}
			fmt.Println("Reinitialization was successful!")
		}
	})
}

// acceptedPath create a strux_pkg directory at the specified path.
// Adding the corresponding entry to the database.
func (ic *InitCommand) acceptedPath() {
	ic.createStruxPkgFolder(func() {
		err := db.CreateDbTable()
		if err != nil {
			panic(err.Error())
		} else {
			struxPkgPath, err := db.GetStruxPkgPathValue()
			if err != nil {
				if err != sql.ErrNoRows {
					panic(err)
				}
			}
			// Column entry strux_pkg_path not found
			if struxPkgPath == "" {
				dirpath := filepath.Join(ic.Path, config.StruxPkgName)
				_, err = db.ExecStruxDbQuery(fmt.Sprintf("INSERT INTO strux VALUES (NULL, '%s', NULL, NULL)", dirpath))
				if err != nil {
					panic(err)
				}
				fmt.Println("Initialization was successful!")
			} else {
				// Column entry strux_pkg_path found
				err := ErrDbStruxPkgPathAlreadyExist{PkgPath: struxPkgPath}
				fmt.Println(err.Error())
			}
		}
	})
}

// updateStructPkgDb the update of the entry in the strux_pkg_path column
// is independent of the creation of the strux_pkg directory.
func (ic *InitCommand) updateStructPkgDb() {
	err := db.CreateDbTable()
	if err != nil {
		panic(err.Error())
	} else {
		// If the path leads to the strux_pkg directory.
		if ic.findStruxPkgInPath(ic.Path) {
			if utils.PathExist(ic.Path) {
				_, err := db.ExecStruxDbQuery(fmt.Sprintf("UPDATE strux SET strux_pkg_path = '%s' WHERE id = 1", ic.Path))
				if err != nil {
					panic(err)
				}
				fmt.Println(fmt.Sprintf("The path %s installed in the database.", ic.Path))
			} else {
				err := ErrCurrentStruxPkgPathNotExist{PkgPath: ic.Path}
				fmt.Println(err.Error())
			}
		} else {
			err := ErrPathMustLeadToStruxPkg{Path: ic.Path}
			fmt.Println(err.Error())
		}
	}
}

// createStruxPkgFolder creates a directory strux_pkg at the specified path.
func (ic *InitCommand) createStruxPkgFolder(onCreate func()) {
	dirpath := filepath.Join(ic.Path, config.StruxPkgName)
	if utils.PathExist(dirpath) {
		err := utils.ErrPathAlreadyExist{Path: dirpath}
		fmt.Println(err.Error())
	} else {
		if err := os.Mkdir(dirpath, os.ModePerm); err != nil {
			panic(err)
		} else {
			onCreate()
		}
	}
}

// findStruxPkgInPath determines whether the path leads to the strux_pkg directory.
func (ic *InitCommand) findStruxPkgInPath(path string) bool {
	var sep rune
	splitFn := func(c rune) bool {
		return c == sep
	}

	if strings.Index(path, "/") != -1 {
		sep = '/'
	}
	if strings.Index(path, "\\") != -1 {
		sep = '\\'
	}

	sPath := strings.FieldsFunc(ic.Path, splitFn)
	val := sPath[len(sPath)-1:][0]
	if strings.Trim(val, string(sep)) == config.StruxPkgName {
		return true
	} else {
		return false
	}
}
