package orm

import (
	"context"
	"database/sql"
	"reflect"
)

var (
	_ db = (*sql.DB)(nil)
	_ db = (*sql.Tx)(nil)
)

type db interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
}

// GetOne execute query and get one result
// query and args may be rewritten. see [RewriteQueryAndArgs] for detail
func GetOne[T any](ctx context.Context, db db, query string, args ...any) (data *T, err error) {
	conf := defaultConfig
	args, opts := parseArgs(args...)
	for _, opt := range opts {
		opt(&conf)
	}

	if conf.rewriteQuery {
		query, args = RewriteQueryAndArgs(query, args)
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
	placeholders := getColumnPlaceholder(&conf, indirect(reflect.ValueOf(&obj)), cols)
	err = rows.Scan(placeholders...)
	if err != nil {
		return nil, err
	}

	return &obj, nil
}

// GetMany execute query and get all result
// query and args may be rewritten. see [RewriteQueryAndArgs] for detail
func GetMany[T any](ctx context.Context, db db, query string, args ...any) (data []*T, err error) {
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

	data = make([]*T, 0)
	for rows.Next() {
		var obj T
		placeholders := getColumnPlaceholder(&conf, indirect(reflect.ValueOf(&obj)), cols)
		err = rows.Scan(placeholders...)
		if err != nil {
			return nil, err
		}
		data = append(data, &obj)
	}

	return data, nil
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
