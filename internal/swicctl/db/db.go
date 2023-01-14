package db

import (
	"github.com/hxia043/nest/pkg/db/file"
	sqlite "github.com/hxia043/nest/pkg/db/sqlite3"
)

const (
	dbworkspace = "/var/lib/swic"
	dbname      = "swic.db"
)

var dbtype = "sqlite"

type DbOption interface {
	Add([]byte) error
	Delete(string) error
	Update([]byte) error
	Get(string) (*Registry, error)
	Show() ([]Registry, error)
}

func InitDB() (DbOption, error) {
	switch dbtype {
	case "file":
		fo := file.NewFileOption(dbworkspace, dbname)
		if err := fo.InitDB(); err != nil {
			return nil, err
		}

		return NewFileOption(fo.DbPath), nil
	case "sqlite":
		so := sqlite.NewSqliteOption(dbname, dbworkspace, DsnSuffix)
		if err := so.InitDB(initCreateTableSQL); err != nil {
			return nil, err
		}

		db, err := so.OpenDB()
		if err != nil {
			return nil, err
		}

		return NewSqliteOption(db), nil
	}

	return nil, nil
}
