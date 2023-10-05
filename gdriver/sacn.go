package gdriver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/d3v-friends/go-pure/fnCases"
	"github.com/d3v-friends/gosql/gctx"
	"reflect"
	"strings"
)

// 맵핑을 위해서 all 로 해야 한다.
func scanAll[T any](rows *sql.Rows) (ls []*T, err error) {
	if !isStruct[T]() {
		err = fmt.Errorf("T is not struct")
		return
	}

	var colNameLs = make([]string, 0)
	if colNameLs, err = rows.Columns(); err != nil {
		return
	}

	ls = make([]*T, 0)
	for rows.Next() {
		var item = new(T)
		var typeOf = reflect.TypeOf(*item)
		var valueOf = reflect.ValueOf(item).Elem()

		var fieldNum = typeOf.NumField()
		var args = make([]any, 0)

		for _, colName := range colNameLs {
			colName = strings.ToLower(colName)

			for i := 0; i < fieldNum; i++ {
				var fieldNm = fnCases.SnakeCase(typeOf.Field(i).Name)
				if colName != fieldNm {
					continue
				}

				args = append(args, valueOf.Field(i).Addr().Interface())
			}
		}

		if len(args) != fieldNum {
			err = fmt.Errorf("struct is not same query result")
			return
		}

		if err = rows.Scan(args...); err != nil {
			return
		}

		ls = append(ls, item)
	}

	return
}

func scanOne[T any](rows *sql.Rows) (res *T, err error) {
	var ls []*T
	if ls, err = scanAll[T](rows); err != nil {
		return
	}

	if len(ls) != 1 {
		err = fmt.Errorf("query is not one result: len=%d", len(ls))
	}

	res = ls[0]
	return
}

func query[T any](
	ctx context.Context,
	query string,
	args ...any,
) (ls []*T, err error) {
	var db = gctx.GetDBP(ctx)
	var rows *sql.Rows
	if rows, err = db.QueryContext(ctx, query, args...); err != nil {
		return
	}
	return scanAll[T](rows)
}

func queryOne[T any](
	ctx context.Context,
	query string,
	args ...any,
) (res *T, err error) {
	var rows *sql.Rows
	var db = gctx.GetDBP(ctx)
	if rows, err = db.QueryContext(ctx, query, args...); err != nil {
		return
	}
	return scanOne[T](rows)
}

func isStruct[T any]() (res bool) {
	return reflect.TypeOf(*new(T)).Kind() == reflect.Struct
}
