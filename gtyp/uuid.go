package gtyp

import (
	"database/sql/driver"
	"fmt"
	"github.com/google/uuid"
	"reflect"
)

type UUID string

func NewUUID() UUID {
	return UUID(uuid.NewString())
}

func (x *UUID) Value() (res driver.Value, _ error) {
	res = string(*x)
	return
}

func (x *UUID) Scan(src any) (err error) {
	switch v := src.(type) {
	case string:
		*x = UUID(v)
		return
	case []byte:
		*x = UUID(v)
		return
	default:
		err = fmt.Errorf("invalid type: src=%s, kind=%s", src, reflect.TypeOf(src).Kind())
		return
	}
}

func (x *UUID) String() string {
	return string(*x)
}
