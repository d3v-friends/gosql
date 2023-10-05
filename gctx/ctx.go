package gctx

import (
	"context"
	"database/sql"
	"fmt"
)

const (
	CtxDB = "CTX_DB"
)

func SetDB(ctx context.Context, db *sql.DB) context.Context {
	ctx = context.WithValue(ctx, CtxDB, db)
	return ctx
}

func GetDB(ctx context.Context) (res *sql.DB, err error) {
	var has bool
	if res, has = ctx.Value(CtxDB).(*sql.DB); !has {
		err = fmt.Errorf("not found *sql.DB in context")
		return
	}
	return
}

func GetDBP(ctx context.Context) (res *sql.DB) {
	var err error
	if res, err = GetDB(ctx); err != nil {
		panic(err)
	}
	return
}
