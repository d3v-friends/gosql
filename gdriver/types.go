package gdriver

import (
	"context"
	"fmt"
	"github.com/d3v-friends/gosql"
)

type Driver interface {
	HasSchema(ctx context.Context, schemaNm gosql.SchemaNm) (has bool, err error)
	CreateSchema(ctx context.Context, schemaNm gosql.SchemaNm) (err error)
	HasTable(ctx context.Context, schemaNm gosql.SchemaNm, modelNm gosql.ModelNm) (has bool, err error)
	CreateTable(ctx context.Context, schemaNm gosql.SchemaNm, modelNm gosql.ModelNm) (err error)
	UpdateColumns(ctx context.Context, schemaNm gosql.SchemaNm, model *gosql.Model) (err error)
	UpdateIndexes(ctx context.Context, schemaNm gosql.SchemaNm, model *gosql.Model) (err error)
}

/* ------------------------------------------------------------------------------------------------------------  */

type count struct {
	Cnt int64
}

/* ------------------------------------------------------------------------------------------------------------  */

type columnInfo struct {
	ColumnName string
	IsNullable string
	DataType   string
	ColumnKey  string
	ColumnType string
}

func (x columnInfo) IsSame(
	model *gosql.Model,
	columnNm gosql.ColumnNm,
	typeMap gosql.TypeMap,
) (res bool, err error) {
	var col = model.Column[columnNm]

	var dType gosql.DType
	if dType, err = typeMap.DType(col); err != nil {
		res = false
		return
	}

	var colType = fmt.Sprintf("%s%s", dType.String(), col.Size())
	var isSameType = colType == x.ColumnType
	var isNullable = (x.IsNullable == "YES") == col.IsPointer()

	res = isSameType && isNullable
	return
}

/* ------------------------------------------------------------------------------------------------------------  */

type indexInfo struct {
	Table        string
	NonUnique    bool
	KeyName      string
	SeqInIndex   int64
	ColumnName   string
	Collation    string
	Cardinality  int64
	SubPart      *int64
	Packed       *string
	Null         string
	IndexType    string
	Comment      string
	IndexComment string
}

type indexInfos []*indexInfo

func (x indexInfos) GetByColumn(columnNm string) (res *indexInfo, err error) {
	for _, info := range x {
		if info.ColumnName == columnNm {
			res = info
			return
		}
	}
	err = fmt.Errorf("not found index info: columnNm=%s", columnNm)
	return
}
