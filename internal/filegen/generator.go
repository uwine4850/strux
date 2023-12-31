package filegen

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"strux/internal/config"
	"time"
)

type FileData struct {
	PkgName     string
	Version     string
	Description string
	InitDate    string
}

type TomlGenerator struct {
	PkgName     string
	Version     string
	Description string
	PkgPath     string
}

// Generate creates a toml configuration file.
func (tg *TomlGenerator) Generate() {
	fileData := FileData{
		PkgName:     tg.PkgName,
		Version:     tg.Version,
		InitDate:    time.Now().Format("02-01-2006"),
		Description: tg.Description}
	file, err := os.Create(filepath.Join(tg.PkgPath, config.ProjectConfName))
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)
	if err := toml.NewEncoder(file).Encode(fileData); err != nil {
		panic(err)
	}
}
