package orm

import (
	"context"
	"database/sql"
	"reflect"
)

func GetOne[T any](ctx context.Context, db *sql.DB, query string, args ...any) (data *T, err error) {
	conf := defaultConfig
	args, opts := parseArgs(args...)
	for _, opt := range opts {
		opt(&conf)
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

func parseArgs(args ...any) (actualArgs []any, opts []callOpt) {
	for _, arg := range args {
		opt, ok := arg.(callOpt)
		if ok {
			opts = append(opts, opt)
		} else {
			actualArgs = append(actualArgs, arg)
		}
	}
	return
}
