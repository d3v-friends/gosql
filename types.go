package gosql

import (
	"fmt"
	"github.com/d3v-friends/go-pure/fnCases"
	"github.com/d3v-friends/go-pure/fnPanic"
	"github.com/gertd/go-pluralize"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

var plu = pluralize.NewClient()

/* ------------------------------------------------------------------------------------------------------------  */

type Version string

/* ------------------------------------------------------------------------------------------------------------  */

type Path string

func NewPath(p string) (path Path) {
	return Path(p)
}

func (x Path) Path() (path string, err error) {
	var ls = filepath.SplitList(x.String())
	if strings.HasPrefix(x.String(), "/") {
		path = filepath.Join(ls...)
		return
	}

	var wd string
	if wd, err = os.Getwd(); err != nil {
		return
	}

	var wdLs = filepath.SplitList(wd)

	wdLs = append(wdLs, ls...)
	path = filepath.Join(wdLs...)
	return
}

func (x Path) String() string {
	return string(x)
}

/* ------------------------------------------------------------------------------------------------------------  */

type Driver string

func (x Driver) String() string {
	return string(x)
}

func (x Driver) IsValid() bool {
	for _, driver := range DriverAll {
		if x == driver {
			return true
		}
	}
	return false
}

const (
	DriverMySQL5 Driver = "mysql5"
)

var DriverAll = []Driver{
	DriverMySQL5,
}

/* ------------------------------------------------------------------------------------------------------------  */

type SDTypeMap map[SType]DType

func (x SDTypeMap) IsValid() (err error) {
	// 기본타입 모두 존재하는지 검사
	for _, sType := range STypeAll {
		var dType, has = x[sType]
		if !has {
			err = fmt.Errorf("no has matched dType with goql type: sType=%s", sType)
			return
		}

		if dType.IsEmpty() {
			err = fmt.Errorf("no has default dType: sType=%s", sType)
			return
		}

	}
	return
}

/* ------------------------------------------------------------------------------------------------------------  */

type SType string
type STypeMap map[string]SType

var sTypeValidate = fnPanic.OnValue(regexp.Compile(`^\*?[a-z]+($|\([0-9]+\)$|\([0-9]+,[0-9]+\)$)`))
var sTypeType = fnPanic.OnValue(regexp.Compile(`[a-z]+`))
var sTypeSize = fnPanic.OnValue(regexp.Compile(`\(([0-9]+(,[0-9]+)?\))`))

const (
	STypesString SType = "string"
	STypeInt     SType = "int"
	STypeFloat   SType = "float"
	STypeBytes   SType = "bytes"
	STypeUUID    SType = "uuid"
	STypeTime    SType = "time"
)

var STypeAll = []SType{
	STypesString,
	STypeInt,
	STypeFloat,
	STypeBytes,
	STypeUUID,
	STypeTime,
}

func (x SType) IsValid() bool {
	return sTypeValidate.MatchString(x.String())
}

func (x SType) getTypeName() SType {
	return SType(sTypeType.FindString(x.String()))
}

func (x SType) IsEqual(v SType) bool {
	return x.getTypeName() == v.getTypeName()
}

func (x SType) IsPointer() bool {
	return strings.HasPrefix(x.String(), "*")
}

func (x SType) Size() string {
	return sTypeSize.FindString(x.String())
}

func (x SType) String() string {
	return string(x)
}

/* ------------------------------------------------------------------------------------------------------------  */

type AType struct {
	SType SType `json:"-" yaml:"-"`
	DType DType `json:"dType" yaml:"dType"`
	GType GType `json:"gType" yaml:"gType"`
}

var defaultTypeMap = TypeMap{
	STypesString: &AType{
		SType: "string",
		DType: "",
		GType: "string",
	},
	STypeInt: &AType{
		SType: "int",
		DType: "",
		GType: "int64",
	},
	STypeFloat: &AType{
		SType: "float",
		DType: "",
		GType: "float",
	},
	STypeBytes: &AType{
		SType: "bytes",
		DType: "",
		GType: "[]byte",
	},
	STypeUUID: &AType{
		SType: "uuid",
		DType: "",
		GType: "github.com/d3v-friends/gosql/gtyp.UUID",
	},
	STypeTime: &AType{
		SType: "time",
		DType: "",
		GType: "time.Time",
	},
}

/* ------------------------------------------------------------------------------------------------------------  */

type DType string

const NilDType DType = ""

var dTypeValidate = fnPanic.OnValue(regexp.Compile(`^\[a-z]+($|\([0-9]+\)$|\([0-9]+,[0-9]+\)$)`))
var dTypeType = fnPanic.OnValue(regexp.Compile(`[a-z]+`))
var dTypeSize = fnPanic.OnValue(regexp.Compile(`\(([0-9]+(,[0-9]+)?\))`))

func (x DType) IsValid() bool {
	return dTypeValidate.MatchString(x.String())
}

func (x DType) IsEmpty() bool {
	return x.String() == ""
}

func (x DType) String() string {
	return string(x)
}

func (x DType) Type() string {
	return dTypeType.FindString(x.String())
}

func (x DType) Size() string {
	return dTypeSize.FindString(x.String())
}

/* ------------------------------------------------------------------------------------------------------------  */

type GType string

func (x GType) GoType() string {
	var ls = strings.Split(x.String(), "/")
	return ls[len(ls)-1]
}

func (x GType) GoImport() string {
	return x.String()
}

func (x GType) String() string {
	return string(x)
}

/* ------------------------------------------------------------------------------------------------------------  */

type TypeMap map[SType]*AType

func (x *TypeMap) IsValid(v SType) (has bool) {
	_, has = (*x)[v.getTypeName()]
	return
}

func (x *TypeMap) DType(v SType) (dType DType, err error) {
	var i, has = (*x)[v.getTypeName()]
	if !has {
		err = fmt.Errorf("not found dType: sType=%s", v)
		return
	}
	dType = i.DType
	return
}

func (x *TypeMap) SetDType(sdTypeMap SDTypeMap) (err error) {
	for sType, dType := range sdTypeMap {
		var aType, has = (*x)[sType]
		if !has {
			err = fmt.Errorf(
				"not found default sType defintion: sType=%s, typeMap=%s",
				sType,
				sdTypeMap,
			)
			return
		}
		aType.DType = dType
	}
	return
}

/* ------------------------------------------------------------------------------------------------------------  */

type ColumnNm string

func (x ColumnNm) String() string {
	return string(x)
}

func (x ColumnNm) SnakeCase() string {
	return fnCases.SnakeCase(x.String())
}

type ColumnMap map[ColumnNm]SType

/* ------------------------------------------------------------------------------------------------------------  */

type Datetime string

const (
	DatetimeCreatedAt Datetime = "createdAt"
	DatetimeUpdatedAt Datetime = "updatedAt"
	DatetimeDeletedAt Datetime = "deletedAt"
)

var DatetimeAll = []Datetime{
	DatetimeCreatedAt,
	DatetimeUpdatedAt,
	DatetimeDeletedAt,
}

func (x Datetime) IsValid() bool {
	for _, datetime := range DatetimeAll {
		if datetime == x {
			return true
		}
	}
	return false
}

type DatetimeList []Datetime

func (x DatetimeList) IsValid() bool {
	for _, datetime := range x {
		if !datetime.IsValid() {
			return false
		}
	}
	return true
}

/* ------------------------------------------------------------------------------------------------------------  */

type Index struct {
	Columns [][]string `json:"columns" yaml:"columns"`
	Unique  bool       `json:"unique" yaml:"unique"`
}

func (x Index) IndexNm() (idxNm string) {
	var ls []string
	for _, column := range x.Columns {
		ls = append(ls, fmt.Sprintf("%s_%s", fnCases.SnakeCase(column[0]), column[1]))
	}

	sort.Strings(ls)

	for _, idx := range ls {
		idxNm = fmt.Sprintf("%s_%s", idxNm, idx)
	}

	idxNm = idxNm[1:]
	return
}

type IndexList []*Index

/* ------------------------------------------------------------------------------------------------------------  */

type Relation struct {
	From         ModelNm      `json:"from" yaml:"from"`
	As           ColumnNm     `json:"as" yaml:"as"`
	LocalField   ColumnNm     `json:"localField" yaml:"localField"`
	ForeignField ColumnNm     `json:"foreignField" yaml:"foreignField"`
	Type         RelationType `json:"type" yaml:"type"`
}

type RelationList []*Relation

type RelationType string

const (
	RelationTypeOne  RelationType = "one"
	RelationTypeMany RelationType = "many"
)

var RelationTypeAll = []RelationType{
	RelationTypeOne,
	RelationTypeMany,
}

func (x RelationType) IsValid() bool {
	for _, relation := range RelationTypeAll {
		if relation == x {
			return true
		}
	}
	return false
}

/* ------------------------------------------------------------------------------------------------------------  */

type Order string

const (
	OrderASC  Order = "asc"
	OrderDESC Order = "desc"
)

var OrderAll = []Order{
	OrderASC,
	OrderDESC,
}

func (x Order) IsValid() bool {
	for _, order := range OrderAll {
		if order == x {
			return true
		}
	}
	return false
}

/* ------------------------------------------------------------------------------------------------------------  */

type ModelNm string

func (x ModelNm) String() string {
	return string(x)
}

func (x ModelNm) SnakeCase() string {
	return fnCases.SnakeCase(plu.Plural(x.String()))
}

func (x ModelNm) SlashCase(schemaNm SchemaNm) string {
	return fmt.Sprintf("%s/%s", schemaNm.SnakeCase(), x.SnakeCase())
}

/* ------------------------------------------------------------------------------------------------------------  */

type SchemaNm string

func (x SchemaNm) String() string {
	return string(x)
}

func (x SchemaNm) SnakeCase() string {
	return fnCases.SnakeCase(x.String())
}
