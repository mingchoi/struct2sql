package struct2sql

import (
	"context"
	"database/sql"
)

type DB struct {
	*sql.DB
}

type Tx struct {
	*sql.Tx
}

type IDB interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	CreateTable(model interface{}) error
	DropTable(model ...interface{}) error
	Insert(model interface{}) error
	Select(model interface{}, condition string, cargs ...interface{}) error
	Update(model interface{}, condition string, cargs ...interface{}) error
	Delete(model interface{}, condition string, cargs ...interface{}) error
}

func Open(driverName, dataSourceName string) (*DB, error) {
	sqlDB, err := sql.Open(driverName, dataSourceName)
	return &DB{sqlDB}, err
}

func (db DB) Begin() (*Tx, error) {
	tx, err := db.DB.Begin()
	return &Tx{tx}, err
}
