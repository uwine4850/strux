package commands

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/uwine4850/strux_api/pkg/uplutils"
	"github.com/uwine4850/strux_api/services/protofiles/baseproto"
	"github.com/uwine4850/strux_api/services/protofiles/pkgproto"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"strux/internal/apiutils"
	"strux/internal/config"
	"strux/internal/db"
	"strux/internal/filegen"
	"strux/utils"
	"syscall"
)

type PkgManagerCommand struct {
	Package    string `short:"pkg" long:"package" block:"1"`
	AddDir     string `short:"ad" long:"add-dir"`
	Rewrite    string `short:"rw" long:"rewrite"`
	Upload     string `short:"upl" long:"upload"`
	Download   string `short:"dwn" long:"download"`
	isUpload   bool
	isAddDir   bool
	isRewrite  bool
	isDownload bool
	pkgName    string
	pkgPath    string
	addPath    string

	downloadUser    string
	downloadVersion string
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
	return []string{p.AddDir, p.Upload, p.Download}
}

func (p *PkgManagerCommand) ExecDownload(user string, version string) []string {
	p.isDownload = true
	p.downloadUser = user
	p.downloadVersion = version
	return []string{}
}

func (p *PkgManagerCommand) ExecUpload() []string {
	p.isUpload = true
	return []string{}
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

// uploadPackage Uploads the packet to the server.
// You must confirm the login and password for the operation to be successful. Accordingly, you must have a registered account.
func (p *PkgManagerCommand) uploadPackage(username string, password string, version string) (*baseproto.BaseResponse, error) {
	pkgName := filepath.Base(p.pkgPath)
	dirInfo, err := uplutils.GetDirsInfo(p.pkgPath, pkgName)
	if err != nil {
		return nil, err
	}
	// dir info to json
	var dirInfoJson []byte
	err = uplutils.UploadDirInfoToJson(dirInfo, &dirInfoJson)
	if err != nil {
		return nil, err
	}

	s, _ := filepath.Split(p.pkgPath)
	pkgUplPath := s[:len(s)-1]
	uplDirInfoFromPaths, err := uplutils.CreateUploadFilePaths(dirInfo, pkgUplPath)
	if err != nil {
		return nil, err
	}
	uploadFilesMap := uplutils.UplFilesToMap(uplDirInfoFromPaths)
	api := apiutils.NewApiForm{
		Method: "POST",
		Url:    "http://localhost:3000/upload-pkg/",
		TextValues: map[string]string{"username": username, "password": password, "version": version,
			"files_info": string(dirInfoJson)},
		FileValues: uploadFilesMap,
	}
	baseResponse, _, err := api.SendForm()
	if err != nil {
		return nil, err
	}
	return baseResponse, nil
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
		for s := range m {
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
		"mkPath":     {mkPath},
		"sourcePath": {sourcePath},
		"mkDir":      {},
		"files":      {},
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
	if p.isUpload {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter username: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		fmt.Print("Enter password(hidden): ")
		bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
		if err != nil {
			panic(err)
		}
		fmt.Print("\n")
		version, err := getPackageInfoVersion(p.pkgName)
		if err != nil {
			panic(err)
		}
		uploadResponse, err := p.uploadPackage(strings.TrimSpace(username), strings.TrimSpace(string(bytePassword)), version)
		if err != nil {
			panic(err)
		}
		fmt.Println(uploadResponse.Message)
		return
	}
	if p.isDownload {
		err := downloadPackage(p.pkgName, p.downloadVersion, p.downloadUser)
		if err != nil {
			panic(err)
		}
		return
	}
}

// downloadPackage Downloads a packet from the server.
// If the package name, version and user are correct, creates the required directory structure and writes the package files there.
func downloadPackage(pkgName string, version string, packageUserOwner string) error {
	downloadPath, err := getDownloadFolderPath(pkgName, version, packageUserOwner)
	if err != nil {
		return err
	}
	api := apiutils.NewApiForm{
		Method:     "GET",
		Url:        "http://localhost:3000/download-package/",
		TextValues: map[string]string{"username": packageUserOwner, "pkgName": pkgName, "version": version},
		FileValues: nil,
	}
	baseResponse, form, err := api.SendForm()
	if err != nil {
		panic(err)
	}
	if baseResponse != nil {
		fmt.Println(baseResponse.Message)
		return nil
	}
	uplDirInfo, err := filesInfoToUploadDirInfo(form.Value["files_info"][0])
	if err != nil {
		panic(err)
	}
	dirTreeMap := make(map[string][]string)
	err = uplutils.CreateDirTree(downloadPath, uplDirInfo, &dirTreeMap)
	if err != nil {
		panic(err)
	}
	var uplFiles []*pkgproto.UploadFile
	err = uplutils.SetUploadFiles(form.File, &uplFiles)
	if err != nil {
		panic(err)
	}
	err = uplutils.CreateFiles(downloadPath, &uplFiles, dirTreeMap)
	if err != nil {
		panic(err)
	}
	return nil
}

// getDownloadFolderPath Creates a folder for downloading the selected package. And returns the path to it.
func getDownloadFolderPath(pkgName string, version string, packageUserOwner string) (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}
	downloadFolder := filepath.Join(currentUser.HomeDir, "StruxDownloads", packageUserOwner, pkgName, version)
	if !utils.PathExist(downloadFolder) {
		err := os.MkdirAll(downloadFolder, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return downloadFolder, nil
}

// filesInfoToUploadDirInfo converts json in string format to a pkgproto.UploadDirInfo structure.
func filesInfoToUploadDirInfo(filesInfo string) (*pkgproto.UploadDirInfo, error) {
	var uplDirInfo *pkgproto.UploadDirInfo
	err := json.Unmarshal([]byte(filesInfo), &uplDirInfo)
	if err != nil {
		return nil, err
	}
	return uplDirInfo, nil
}

// getPackageInfoVersion read the project.toml file and return the current version.
func getPackageInfoVersion(pkgName string) (string, error) {
	pathPkg, err := db.GetStruxPkgPathValue()
	if err != nil {
		return "", err
	}
	projectTomlPath := filepath.Join(pathPkg, pkgName, config.ProjectConfName)
	if utils.PathExist(projectTomlPath) {
		var fd filegen.FileData
		file, err := os.ReadFile(projectTomlPath)
		if err != nil {
			return "", err
		}
		// get/set data
		_, err = toml.Decode(string(file), &fd)
		if err != nil {
			return "", err
		}
		return fd.Version, nil
	}
	return "", errors.New(fmt.Sprintf("Path %s from project.toml not exist", projectTomlPath))
}
