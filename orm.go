package orm

import (
	"context"
	"database/sql"
	"fmt"
)

var (
	ErrConcurrencyUpdate = fmt.Errorf("concurrency update")
)

var (
	_ DB = (*sql.DB)(nil)
	_ DB = (*sql.Tx)(nil)
)

type DB interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

// GetOne execute query and get one result
// query and args may be rewritten. see [RewriteQueryAndArgs] for detail
func GetOne[T any](ctx context.Context, db DB, query string, args ...any) (*T, error) {
	conf := defaultConfig
	args, opts := parseArgs(args...)
	for _, opt := range opts {
		opt(&conf)
	}

	if conf.rewriteQuery {
		query, args = RewriteQueryAndArgs(query, args...)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var obj T
	dest := getColumnDest(&conf, &obj, cols)
	err = rows.Scan(dest...)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

// GetMany execute query and get all result
// query and args may be rewritten. see [RewriteQueryAndArgs] for detail
func GetMany[T any](ctx context.Context, db DB, query string, args ...any) ([]*T, error) {
	conf := defaultConfig
	args, opts := parseArgs(args...)
	for _, opt := range opts {
		opt(&conf)
	}

	if conf.rewriteQuery {
		query, args = RewriteQueryAndArgs(query, args...)
	}

	rows, err := db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	data := make([]*T, 0)
	for rows.Next() {
		var obj T
		dest := getColumnDest(&conf, &obj, cols)
		err = rows.Scan(dest...)
		if err != nil {
			return nil, err
		}
		data = append(data, &obj)
	}

	return data, nil
}

func InsertOne(ctx context.Context, db DB, tableName string, data any, opts ...func(*config)) error {
	conf := defaultConfig
	for _, opt := range opts {
		opt(&conf)
	}

	values := tagParser.Parse(conf.tagName, data)
	insertColumns, autoIncrementColumn, args := parseInsertColumnsAndArgs(values)
	query := generateInsertSQL(tableName, insertColumns, 1)
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	if autoIncrementColumn != "" && values.Get(autoIncrementColumn).CanSet() {
		lastInsertId, err := result.LastInsertId()
		if err != nil {
			return err
		}
		values.Get(autoIncrementColumn).Set(lastInsertId)
	}

	return nil
}

func InsertMany[T any](ctx context.Context, db DB, tableName string, data []T, opts ...func(*config)) error {
	if len(data) == 0 {
		return nil
	}

	conf := defaultConfig
	for _, opt := range opts {
		opt(&conf)
	}

	values := tagParser.Parse(conf.tagName, data[0])
	insertColumns, _, _ := parseInsertColumnsAndArgs(values)

	batchSize := min(conf.batchSize, len(data))
	query := generateInsertSQL(tableName, insertColumns, batchSize)
	args := make([]any, 0, len(insertColumns)*batchSize)

	for i := 0; i < len(data); i += batchSize {
		args = args[:0]
		end := min(i+batchSize, len(data))
		batch := data[i:end]
		if len(batch) < batchSize {
			query = generateInsertSQL(tableName, insertColumns, len(batch))
		}
		for _, item := range batch {
			itemValues := tagParser.Parse(conf.tagName, item)
			for _, col := range insertColumns {
				args = append(args, itemValues.Get(col).Interface())
			}
		}
		_, err := db.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
	}

	return nil
}

func UpdateOne(ctx context.Context, db DB, tableName string, data any, opts ...func(*config)) error {
	conf := defaultConfig
	for _, opt := range opts {
		opt(&conf)
	}

	values := tagParser.Parse(conf.tagName, data)
	updateColumns, whereColumns, updateArgs, whereArgs, versionValue := parseUpdateColumnsAndArgs(&conf, values)
	query := generateUpdateSQL(tableName, updateColumns, whereColumns)
	args := append(updateArgs, whereArgs...)
	result, err := db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if versionValue != nil {
		if rowsAffected == 0 {
			return ErrConcurrencyUpdate
		}
		if versionValue.CanSet() {
			versionValue.Set(versionValue.Value().Int() + 1)
		}
	}

	return nil
}

func parseArgs(args ...any) (actualArgs []any, opts []func(*config)) {
	for _, arg := range args {
		opt, ok := arg.(func(*config))
		if ok {
			opts = append(opts, opt)
		} else {
			actualArgs = append(actualArgs, arg)
		}
	}
	return
}
