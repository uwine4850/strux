package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"strux/internal/db"
	"strux/utils"
)

type PkgManagerCommand struct {
	Package   string `short:"pkg" long:"package" block:"1"`
	AddDir    string `short:"ad" long:"adddir"`
	Rewrite   string `short:"rw" long:"rewrite"`
	isAddDir  bool
	isRewrite bool
	pkgName   string
	pkgPath   string
	addPath   string
}

func (p *PkgManagerCommand) ExecPackage(pkgName string) []string {
	struxPkgPath, err := db.GetStruxPkgPathValue()
	if err != nil {
		panic(err)
	}
	pkgPath := filepath.Join(struxPkgPath, pkgName)
	if utils.PathExist(pkgPath) {
		p.pkgName = pkgName
		p.pkgPath = pkgPath
	} else {
		err := &ErrPkgOrConfigNotFound{PkgName: pkgName}
		fmt.Println(err.Error())
	}
	return []string{p.AddDir}
}

func (p *PkgManagerCommand) ExecAddDir(addPath string) []string {
	p.isAddDir = true
	if !utils.PathExist(addPath) {
		err := &utils.ErrPathNotExist{Path: addPath}
		panic(err)
	} else {
		p.addPath = addPath
	}
	return []string{p.Rewrite}
}

func (p *PkgManagerCommand) ExecRewrite() []string {
	p.isRewrite = true
	return []string{}
}

// addDir creates directories and their files in the selected package.
func (p *PkgManagerCommand) addDir() {
	stat, err := os.Stat(p.addPath)
	if err != nil {
		panic(err)
	}
	var dirList []map[string][]string
	if stat.IsDir() {
		p.setAddDirList(p.pkgPath, p.addPath, &dirList)
	} else {
		err := &utils.ErrThisIsNotADir{DirPath: p.addPath}
		panic(err)
	}
	var mkPath string
	var sourcePath string
	var mkDir string
	var files []string
	for _, m := range dirList {
		for s, _ := range m {
			switch s {
			// make files this
			case "mkPath":
				mkPath = m[s][0]
			// folder from source files
			case "sourcePath":
				sourcePath = m[s][0]
			// mk new dir
			case "mkDir":
				mkDir = m[s][0]
			// mk files
			case "files":
				files = m[s]
			}
		}

		// add dirs
		mkPath := filepath.Join(mkPath, mkDir)
		if !utils.PathExist(mkPath) {
			err := os.Mkdir(mkPath, os.ModePerm)
			if err != nil {
				panic(err)
			}
			if !utils.PathExist(mkPath) {
				err := &utils.ErrDirNotCreated{DirPath: mkPath}
				panic(err)
			}
		}

		// add files
		if utils.PathExist(mkPath) {
			for i := 0; i < len(files); i++ {
				filePath := filepath.Join(sourcePath, files[i])
				if !utils.PathExist(filePath) {
					err := &utils.ErrPathNotExist{Path: filePath}
					panic(err)
				}
				fileName := filepath.Base(filePath)
				newFilePath := filepath.Join(mkPath, fileName)
				// copying files
				err := utils.CopyFile(filePath, newFilePath, p.isRewrite)
				if err != nil {
					panic(err)
				}
			}
		}
	}
}

// setAddDirList sets a specially formatted list of files and paths to the corresponding directories.
func (p *PkgManagerCommand) setAddDirList(mkPath string, sourcePath string, dirList *[]map[string][]string) {
	m := map[string][]string{
		"mkPath":     []string{mkPath},
		"sourcePath": []string{sourcePath},
		"mkDir":      []string{},
		"files":      []string{},
	}
	dirFiles, err := os.ReadDir(sourcePath)
	if err != nil {
		panic(err)
	}
	// add mkDir
	sourcePath = strings.ReplaceAll(sourcePath, "\\", "/")
	mkDir := strings.Split(sourcePath, "/")[len(strings.Split(sourcePath, "/"))-1:][0]
	m["mkDir"] = append(m["mkDir"], mkDir)

	// collect dirs from current path
	var dirs []string
	for i := 0; i < len(dirFiles); i++ {
		if dirFiles[i].IsDir() {
			dirs = append(dirs, dirFiles[i].Name())
		} else {
			m["files"] = append(m["files"], dirFiles[i].Name())
		}
	}
	// set all data
	*dirList = append(*dirList, m)

	// processing next dir
	for i := 0; i < len(dirs); i++ {
		newSourcePath := filepath.Join(sourcePath, dirs[i])
		newMkPath := filepath.Join(mkPath, mkDir)
		if utils.PathExist(newSourcePath) {
			p.setAddDirList(newMkPath, newSourcePath, dirList)
		} else {
			err := &utils.ErrPathNotExist{Path: newSourcePath}
			panic(err)
		}
	}
}

func (p *PkgManagerCommand) OnFinish() {
	if p.isAddDir {
		p.addDir()
	}
}
