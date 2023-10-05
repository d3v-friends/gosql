package gdriver

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/d3v-friends/go-pure/fnCases"
	"github.com/d3v-friends/go-pure/fnLogger"
	"github.com/d3v-friends/go-pure/fnMatch"
	"github.com/d3v-friends/gosql"
	"github.com/d3v-friends/gosql/gctx"
	_ "github.com/go-sql-driver/mysql"
	"sort"
	"time"
)

type MySQL5 struct {
	typeMap gosql.TypeMap
}

var mysql5SdType = gosql.SDTypeMap{
	gosql.STypesString: "varchar",
	gosql.STypeInt:     "bigint",
	gosql.STypeFloat:   "double",
	gosql.STypeBytes:   "blob",
	gosql.STypeUUID:    "char",
	gosql.STypeTime:    "datetime",
}

type IConn struct {
	Username string
	Password string
	Host     string
}

func (x IConn) mySql5Host() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s)/?charset=utf8mb4&time_zone=utc&parseTime=true",
		x.Username,
		x.Password,
		x.Host,
	)
}

func NewMySQL5(
	typeMap gosql.TypeMap,
	i *IConn,
) (
	driver *MySQL5,
	db *sql.DB,
	err error,
) {
	driver = &MySQL5{
		typeMap: typeMap,
	}

	if err = driver.typeMap.SetDType(mysql5SdType); err != nil {
		return
	}

	if db, err = sql.Open("mysql", i.mySql5Host()); err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		return
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return
}

func (x *MySQL5) HasSchema(ctx context.Context, schemaNm gosql.SchemaNm) (has bool, err error) {
	var logger = x.logger(ctx)
	var execSql = fmt.Sprintf(mysql5CountSchema, schemaNm.SnakeCase())

	logger.Trace(execSql)

	var co = &count{}
	if co, err = queryOne[count](ctx, execSql); err != nil {
		return
	}

	has = 0 < co.Cnt
	return
}

func (x *MySQL5) CreateSchema(ctx context.Context, schemaNm gosql.SchemaNm) (err error) {
	var logger = x.logger(ctx)
	var execSql = fmt.Sprintf(mysql5CreateSchema, schemaNm.SnakeCase())
	logger.Trace(execSql)

	var db = gctx.GetDBP(ctx)
	if _, err = db.ExecContext(ctx, execSql); err != nil {
		return
	}
	return
}

func (x *MySQL5) HasTable(
	ctx context.Context,
	schemaNm gosql.SchemaNm,
	modelNm gosql.ModelNm,
) (has bool, err error) {
	var logger = x.logger(ctx)
	var execSql = fmt.Sprintf(mysql5CountTable, modelNm.SlashCase(schemaNm))
	logger.Trace(execSql)

	var co = &count{}

	if co, err = queryOne[count](ctx, execSql); err != nil {
		return
	}

	has = 0 < co.Cnt
	return
}

func (x *MySQL5) CreateTable(
	ctx context.Context,
	schemaNm gosql.SchemaNm,
	modelNm gosql.ModelNm,
) (err error) {
	var logger = x.logger(ctx)
	var execSql = fmt.Sprintf(
		mysql5CreateTable,
		schemaNm.SnakeCase(),
		modelNm.SnakeCase(),
	)

	logger.
		WithFields(fnLogger.Fields{
			"run": "createTable",
		}).
		Trace(execSql)

	var db = gctx.GetDBP(ctx)
	if _, err = db.ExecContext(ctx, execSql); err != nil {
		return
	}
	return
}

func (x *MySQL5) columnInfos(
	ctx context.Context,
	schemaNm gosql.SchemaNm,
	modelNm gosql.ModelNm,
) (ls []*columnInfo, err error) {
	var logger = x.logger(ctx)
	var execSql = fmt.Sprintf(
		mysql5ReadColumnInfo,
		modelNm.SnakeCase(),
		schemaNm.SnakeCase(),
	)

	logger.
		WithFields(fnLogger.Fields{
			"run": "getColumnInfos",
		}).
		Trace(execSql)

	if ls, err = query[columnInfo](ctx, execSql); err != nil {
		return
	}

	return
}

func (x *MySQL5) UpdateColumns(
	ctx context.Context,
	schemaNm gosql.SchemaNm,
	model *gosql.Model,
) (err error) {
	var columnInfos []*columnInfo
	if columnInfos, err = x.columnInfos(ctx, schemaNm, model.ModelNm); err != nil {
		return
	}

	for columnNm := range model.Column {
		var colInfo *columnInfo
		if colInfo, err = fnMatch.Get(columnInfos, func(v *columnInfo) bool {
			return v.ColumnName == columnNm.SnakeCase()
		}); err != nil {
			err = x.createColumn(ctx, model, columnNm)
		} else {
			var isSame bool
			if isSame, err = colInfo.IsSame(model, columnNm, x.typeMap); err != nil {
				return
			}

			if !isSame {
				err = x.alterColumn(ctx, model, columnNm)
			}
		}

		if err != nil {
			return
		}
	}

	return
}

func (x *MySQL5) createColumn(
	ctx context.Context,
	model *gosql.Model,
	columnNm gosql.ColumnNm,
) (err error) {
	var col = model.Column[columnNm]
	var dType gosql.DType
	if dType, err = x.typeMap.DType(col); err != nil {
		return
	}

	var execSql = fmt.Sprintf(
		mysql5CreateColumn,
		model.ModelNm.SnakeCase(),
		columnNm.SnakeCase(),
		dType.Type(),
	)

	if col.Size() != "" {
		execSql = fmt.Sprintf("%s%s", execSql, col.Size())
	}

	if !col.IsPointer() {
		execSql = fmt.Sprintf("%s not null", execSql)
	}

	var logger = x.logger(ctx).WithFields(fnLogger.Fields{
		"run": "create column",
	})
	logger.Trace(execSql)

	var db = gctx.GetDBP(ctx)
	if _, err = db.ExecContext(ctx, execSql); err != nil {
		return
	}

	return
}

func (x *MySQL5) alterColumn(
	ctx context.Context,
	model *gosql.Model,
	columnNm gosql.ColumnNm,
) (err error) {
	var col = model.Column[columnNm]

	var dType gosql.DType
	if dType, err = x.typeMap.DType(col); err != nil {
		return
	}

	var execSql = fmt.Sprintf(
		mysql5AlterColumn,
		model.ModelNm.SnakeCase(),
		columnNm.SnakeCase(),
		dType.Type(),
	)

	var size = col.Size()
	if size != "" {
		execSql = fmt.Sprintf("%s%s", execSql, size)
	}

	if !col.IsPointer() {
		execSql = fmt.Sprintf("%s not null", execSql)
	}

	var logger = x.logger(ctx).WithFields(fnLogger.Fields{
		"run": "alter column",
	})
	logger.Trace(execSql)

	var db = gctx.GetDBP(ctx)
	if _, err = db.ExecContext(ctx, execSql); err != nil {
		return
	}

	return
}

func (x *MySQL5) UpdateIndexes(
	ctx context.Context,
	schemaNm gosql.SchemaNm,
	model *gosql.Model,
) (err error) {
	for _, index := range model.Index {
		var idxNm = index.IndexNm()
		var execSql = fmt.Sprintf(mysql5ReadIndexes, schemaNm.SnakeCase(), model.ModelNm.SnakeCase(), idxNm)
		var ls []*indexInfo

		if ls, err = query[indexInfo](ctx, execSql); err != nil {
			return
		}

		if len(ls) == 0 {
			if err = x.createIndex(ctx, schemaNm, model, index); err != nil {
				return
			}
			continue
		}

		if len(index.Columns) != len(ls) {
			if err = x.dropAndCreateIndex(ctx, schemaNm, model, index); err != nil {
				return
			}
			continue
		}

		// column 모두 있나 확인
		for _, indexDef := range index.Columns {
			var columnNm = indexDef[0]

			var has = fnMatch.Has(ls, func(v *indexInfo) bool {
				return v.ColumnName == fnCases.SnakeCase(columnNm)
			})

			if !has {
				if err = x.dropAndCreateIndex(ctx, schemaNm, model, index); err != nil {
					return
				}
				break
			}
		}
	}

	return
}

func (x *MySQL5) dropAndCreateIndex(
	ctx context.Context,
	schemaNm gosql.SchemaNm,
	model *gosql.Model,
	index *gosql.Index,
) (err error) {
	var execSql = fmt.Sprintf(mysql5DropIndex, schemaNm.SnakeCase(), model.ModelNm.SnakeCase(), index.IndexNm())
	var db = gctx.GetDBP(ctx)

	var logger = x.logger(ctx)
	logger.
		WithFields(fnLogger.Fields{
			"run": "drop index",
		}).
		Trace(execSql)

	if _, err = db.ExecContext(ctx, execSql); err != nil {
		return
	}

	return x.createIndex(ctx, schemaNm, model, index)
}

func (x *MySQL5) createIndex(
	ctx context.Context,
	schemaNm gosql.SchemaNm,
	model *gosql.Model,
	index *gosql.Index,
) (err error) {
	var format = mysql5CreateIndex
	if index.Unique {
		format = mysql5CreateUniqueIndex
	}

	var idxDef = parseIndexValue(index)

	var execSql = fmt.Sprintf(
		format,
		index.IndexNm(),
		schemaNm.SnakeCase(),
		model.ModelNm.SnakeCase(),
		idxDef,
	)

	var logger = x.logger(ctx)
	logger.
		WithFields(fnLogger.Fields{
			"run": "create index",
		}).
		Trace(execSql)

	var db = gctx.GetDBP(ctx)
	if _, err = db.ExecContext(ctx, execSql); err != nil {
		return
	}

	return
}

func (x *MySQL5) logger(ctx context.Context) (res fnLogger.IfLogger) {
	res = fnLogger.Get(ctx, &fnLogger.DummyLogger{})
	return
}

func parseIndexValue(idx *gosql.Index) (res string) {
	var ls = make([]string, 0)
	for _, column := range idx.Columns {
		ls = append(ls, fmt.Sprintf("%s %s", fnCases.SnakeCase(column[0]), column[1]))
	}

	sort.Strings(ls)

	for _, l := range ls {
		res = fmt.Sprintf("%s, %s", res, l)
	}

	res = res[1:]
	return
}
