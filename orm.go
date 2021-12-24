package orm

import (
	"context"
	"database/sql"
	"reflect"
)

type DB struct {
	*sql.DB
}

func Open(driverName, dataSourceName string) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DB{
		DB: db,
	}, err
}

func (db *DB) GetOne(ctx context.Context, obj interface{}, query string, args ...interface{}) (ok bool, err error) {
	objT := reflect.TypeOf(obj)
	if !isStructPtr(objT) {
		panic("todo")
	}
	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if !rows.Next() {
		return false, nil
	}

	cols, err := rows.Columns()
	if err != nil {
		return false, err
	}
	placeholders := getColumnPlaceholder(reflect.ValueOf(obj).Elem(), cols)

	err = rows.Scan(placeholders...)
	if err != nil {
		return false, err
	}

	return true, nil
}

// ============================================

type empty struct{}

func (empty) Scan(src interface{}) error {
	return nil
}

func isStructPtr(t reflect.Type) bool {
	if t.Kind() != reflect.Ptr {
		return false
	}
	if t.Elem().Kind() != reflect.Struct {
		return false
	}
	return true
}
