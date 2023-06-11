package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"strux/internal/config"
)

type Database struct {
	DbPath string
}

func (d *Database) Open() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", d.DbPath)
	if err != nil {
		return nil, err
	}
	return db, err
}

// GetStruxPkgPathValue returns the current path to the strux_pkg directory.
func GetStruxPkgPathValue() (string, error) {
	d := Database{DbPath: config.DbName}
	db, err := d.Open()
	if err != nil {
		return "", err
	}
	row := db.QueryRow("SELECT strux_pkg_path FROM strux WHERE id=?", 1)
	struxPkgPath := ""
	if err = row.Scan(&struxPkgPath); err == sql.ErrNoRows {
		return "", err
	}
	return struxPkgPath, err
}

// CreateDbTable creating a strux table.
func CreateDbTable() error {
	d := Database{DbPath: config.DbName}
	db, err := d.Open()
	if err != nil {
		return err
	}
	_, err = db.Exec(config.StruxTableSql)
	if err != nil {
		return err
	}
	return nil
}

// ExecStruxDbQuery execute a query against a strux table.
func ExecStruxDbQuery(query string) (sql.Result, error) {
	d := Database{DbPath: config.DbName}
	db, err := d.Open()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic(err)
		}
	}(db)
	if err != nil {
		panic(err)
	}
	res, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	return res, err
}
