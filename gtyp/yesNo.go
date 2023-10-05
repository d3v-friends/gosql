package gtyp

import (
	"database/sql/driver"
	"fmt"
)

type YesNo string

func (x *YesNo) Scan(src any) (err error) {
	switch v := src.(type) {
	case string:
		*x = YesNo(v)
	case []byte:
		*x = YesNo(v)
	default:
		err = fmt.Errorf("invalid YesNo value: src=%s", src)
	}
	return
}

func (x *YesNo) Value() (driver.Value, error) {
	return string(*x), nil
}

func (x *YesNo) IsValid() bool {
	for _, no := range YesNoAll {
		if *x == no {
			return true
		}
	}
	return false
}

const (
	YesNoYes YesNo = "YES"
	YesNoNo  YesNo = "NO"
)

var YesNoAll = []YesNo{
	YesNoYes,
	YesNoNo,
}
