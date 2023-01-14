package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type TableSQL string

type SqliteOption struct {
	Name      string
	Workspace string
	DsnSuffix string
}

func (so *SqliteOption) CreateWorkspace() (bool, error) {
	_, err := os.Stat(so.Workspace)
	if err == nil {
		return true, nil
	}

	if err := os.MkdirAll(so.Workspace, os.ModePerm); err != nil {
		return false, err
	}

	return true, nil
}

func (so *SqliteOption) InitDB(tableCreateSQLs []TableSQL) error {
	isCreate, err := so.CreateWorkspace()
	if !isCreate {
		return err
	}

	db, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s%s", so.Workspace, so.Name, so.DsnSuffix))
	if err != nil {
		return err
	}

	defer db.Close()

	for _, tableCreateSQL := range tableCreateSQLs {
		if _, err := db.Exec(string(tableCreateSQL)); err != nil {
			return err
		}
	}

	return nil
}

func (so *SqliteOption) OpenDB() (*sql.DB, error) {
	dbconnection, err := sql.Open("sqlite3", fmt.Sprintf("%s/%s%s", so.Workspace, so.Name, so.DsnSuffix))
	if err != nil {
		log.Println(err)
	}

	return dbconnection, err
}

func NewSqliteOption(name, workspace, dsnSuffix string) *SqliteOption {
	return &SqliteOption{
		Name:      name,
		Workspace: workspace,
		DsnSuffix: dsnSuffix,
	}
}
