package gctl

import (
	"context"
	"database/sql"
	"github.com/d3v-friends/gosql"
	"github.com/d3v-friends/gosql/gctx"
	"github.com/d3v-friends/gosql/gdriver"
)

func Migrate(
	ctx context.Context,
	db *sql.DB,
	driver gdriver.Driver,
	cfg *gosql.Config,
) (err error) {
	var schemaNm = cfg.Schema
	ctx = gctx.SetDB(ctx, db)

	// schema check
	var has bool
	if has, err = driver.HasSchema(ctx, schemaNm); err != nil {
		return
	}

	if !has {
		if err = driver.CreateSchema(ctx, schemaNm); err != nil {
			return
		}
	}

	// table check
	for _, model := range cfg.Model {
		if has, err = driver.HasTable(ctx, schemaNm, model.ModelNm); err != nil {
			return
		}

		if !has {
			if err = driver.CreateTable(ctx, schemaNm, model.ModelNm); err != nil {
				return
			}
		}
	}

	// column updates
	for _, model := range cfg.Model {
		if err = driver.UpdateColumns(ctx, schemaNm, model); err != nil {
			return
		}
	}

	// index updates
	for _, model := range cfg.Model {
		if err = driver.UpdateIndexes(ctx, schemaNm, model); err != nil {
			return
		}
	}

	return
}
