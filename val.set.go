package gosql

import (
	"fmt"
	"os"
)

func (x *Config) valPath() (err error) {
	var path string
	if path, err = x.Out.Path(); err != nil {
		return
	}

	if err = os.MkdirAll(path, os.ModePerm); err != nil {
		return
	}

	return
}

func (x *Config) valDriver() (err error) {
	if x.Driver.IsValid() {
		return
	}
	err = fmt.Errorf("unsopported driver: driver=%s", x.Driver)
	return
}

func (x *Config) valType() (err error) {
	for sType, aType := range x.Type {
		aType.SType = sType.getTypeName()
	}

	var has bool
	for sType, aType := range defaultTypeMap {
		if _, has = x.Type[sType]; has {
			continue
		}
		x.Type[sType] = aType
	}

	// sType 형태검사
	for sType := range x.Type {
		if !sType.IsValid() {
			err = fmt.Errorf("invalid sType: sType=%s", sType)
			return
		}
	}

	// 기본타입 추가하기
	for sType, aType := range defaultTypeMap {
		x.Type[sType] = aType
	}

	return
}

func (x *Config) valColumnType() (err error) {
	var errFmt = "invalid column: modelNm=%s, columnNm=%s, sType=%s"
	for modelNm, model := range x.Model {
		for columnNm, sType := range model.Column {
			if !x.Type.IsValid(sType) {
				err = fmt.Errorf(errFmt, modelNm, columnNm, sType)
				return
			}
		}
	}
	return
}

func (x *Config) setModelNm() (err error) {
	for modelNm, model := range x.Model {
		model.ModelNm = modelNm
	}
	return
}

func (x *Config) valDatetime() (err error) {
	for modelNm, model := range x.Model {
		if model.Datetime.IsValid() {
			continue
		}
		err = fmt.Errorf("invalid datetime: modelNm=%s", modelNm)
	}
	return
}

func (x *Config) valIndex() (err error) {
	for modelNm, model := range x.Model {
		for idx, index := range model.Index {
			for _, idxDef := range index.Columns {
				if len(idxDef) != 2 {
					err = fmt.Errorf("invalid index: modelNm=%s, indexId=%d, index=%s", modelNm, idx, index.Columns)
					return
				}

				var has bool
				var columnNm = ColumnNm(idxDef[0])
				if _, has = model.Column[columnNm]; !has {
					err = fmt.Errorf(
						"not found model index column: modelNm=%s, indexIdx=%d, indexColumnNm=%s",
						modelNm,
						idx,
						columnNm,
					)
					return
				}

				var order = Order(idxDef[1])
				if !order.IsValid() {
					err = fmt.Errorf(
						"invalid order: modelNm=%s, indexIdx=%d, order=%s",
						modelNm,
						idx,
						order,
					)
					return
				}
			}
		}
	}

	return
}

func (x *Config) valRelation() (err error) {
	for modelNm, model := range x.Model {
		for relIdx, relation := range model.Relation {
			var has bool

			// from
			if _, has = x.Model[relation.From]; !has {
				err = fmt.Errorf(
					"not found relatoin model: modelNm=%s, relationIdx=%d, fromModelNm=%s",
					modelNm,
					relIdx,
					relation.From,
				)
				return
			}

			// as
			if _, has = model.Column[relation.As]; has {
				err = fmt.Errorf(
					"already has same name column: modelNm=%s, relationIdx=%d, relationColumnNm=%s",
					modelNm,
					relIdx,
					relation.As,
				)
				return
			}

			// localField
			if _, has = model.Column[relation.LocalField]; !has {
				err = fmt.Errorf(
					"not found local field: modelNm=%s, relatoinIdx= %d, localField=%s",
					modelNm,
					relIdx,
					relation.LocalField,
				)
				return
			}

			// foreignField
			var fromModel *Model
			if fromModel, has = x.Model[relation.From]; !has {
				err = fmt.Errorf(
					"not found foreign model: modelNm=%s, relationIdx=%d, fromModel=%s",
					modelNm,
					relIdx,
					relation.From,
				)
				return
			}

			if _, has = fromModel.Column[relation.ForeignField]; !has {
				err = fmt.Errorf(
					"not found foreign model: modelNm=%s, relationIdx=%d, fromModel=%s",
					modelNm,
					relIdx,
					relation.From,
				)
				return
			}

		}
	}

	return
}
