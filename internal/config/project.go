package config

const StruxPkgName = "strux_pkg"
const DbName = "strux.sqlite"
const ProjectConfName = "project.toml"
const StruxTableSql = `
	CREATE TABLE IF NOT EXISTS strux (
	    id INTEGER NOT NULL PRIMARY KEY,
	    strux_pkg_path VARCHAR(200) NOT NULL,
	    login VARCHAR(200) NULL,
	    password VARCHAR(200) NULL
	)
`
