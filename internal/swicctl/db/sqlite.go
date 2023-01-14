package db

import (
	"database/sql"
	"encoding/json"
	"sync"

	sqlite "github.com/hxia043/nest/pkg/db/sqlite3"
)

const (
	DsnSuffix              = "?_fk=true&_busy_timeout=30000"
	registryTableName      = "registry"
	createRegistryTableSQL = sqlite.TableSQL(`create table if not exists Registry(
		name text not null primary key,
		username text not null, 
		password text not null,
		cert text not null, 
		registry text not null);`,
	)
)

var initCreateTableSQL = []sqlite.TableSQL{createRegistryTableSQL}

type SqliteOption struct {
	DB        *sql.DB
	DbRwMutex sync.RWMutex
	Registry
}

func (o *SqliteOption) Add(option []byte) error {
	if err := json.Unmarshal(option, &o.Registry); err != nil {
		return err
	}

	sqlCmd := "insert into Registry (name,username,password,cert,registry) values (?,?,?,?,?)"
	o.DbRwMutex.Lock()
	defer o.DbRwMutex.Unlock()

	if _, err := o.DB.Exec(sqlCmd, o.Name, o.Username, o.Password, o.CaCert, o.Registry.Registry); err != nil {
		return err
	}

	return nil
}

func (o *SqliteOption) Get(name string) (*Registry, error) {
	sqlCmd := "select name, username, password, registry, cert from Registry where name = ?"
	o.DbRwMutex.Lock()
	defer o.DbRwMutex.Unlock()

	if err := o.DB.QueryRow(sqlCmd, name).Scan(&o.Name, &o.Username, &o.Password, &o.Registry.Registry, &o.CaCert); err != nil {
		return nil, err
	}

	return &o.Registry, nil
}

func (o *SqliteOption) Show() ([]Registry, error) {
	sqlCmd := "select name, username, registry from Registry"
	o.DbRwMutex.Lock()
	defer o.DbRwMutex.Unlock()

	rows, err := o.DB.Query(sqlCmd)
	if err != nil {
		return nil, err
	}

	rs := []Registry{}
	for rows.Next() {
		var r Registry
		if err := rows.Scan(&r.Name, &r.Username, &r.Registry); err != nil {
			return nil, err
		}

		rs = append(rs, r)
	}

	return rs, nil
}

func (o *SqliteOption) Update(option []byte) error {
	if err := json.Unmarshal(option, &o.Registry); err != nil {
		return err
	}

	sqlCmd := "update Registry set username=?, password=?, cert=?, registry=? where name=?"
	o.DbRwMutex.Lock()
	defer o.DbRwMutex.Unlock()

	if _, err := o.DB.Exec(sqlCmd, o.Username, o.Password, o.CaCert, o.Registry.Registry, o.Name); err != nil {
		return err
	}

	return nil
}

func (o *SqliteOption) validateRegistryRow(name string) error {
	sqlCmd := "select name from Registry where name = ?"
	o.DbRwMutex.Lock()
	defer o.DbRwMutex.Unlock()

	var rowname string
	if err := o.DB.QueryRow(sqlCmd, name).Scan(&rowname); err != nil {
		return err
	}

	return nil
}

func (o *SqliteOption) Delete(name string) error {
	if err := o.validateRegistryRow(name); err != nil {
		return err
	}

	sqlCmd := "delete from Registry where name = ?"
	o.DbRwMutex.Lock()
	defer o.DbRwMutex.Unlock()

	if _, err := o.DB.Exec(sqlCmd, name); err != nil {
		return err
	}

	return nil
}

func NewSqliteOption(db *sql.DB) *SqliteOption {
	return &SqliteOption{DB: db}
}
