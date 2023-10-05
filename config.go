package gosql

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Version Version  `json:"version" yaml:"version"`
	Out     Path     `json:"out" yaml:"out"`
	Driver  Driver   `json:"driver" yaml:"driver"`
	Schema  SchemaNm `json:"schema" yaml:"schema"`
	Type    TypeMap  `json:"type" yaml:"type"`
	Model   ModelMap `json:"model" yaml:"model"`
}

type Model struct {
	ModelNm  ModelNm      `json:"-" yaml:"-"`
	Column   ColumnMap    `json:"column" yaml:"column"`
	Datetime DatetimeList `json:"datetime" yaml:"datetime"`
	Index    IndexList    `json:"index" json:"index"`
	Relation RelationList `json:"relation" json:"relation"`
}

type ModelMap map[ModelNm]*Model

func Read(fp Path) (cfg *Config, err error) {
	var path string
	if path, err = fp.Path(); err != nil {
		return
	}

	var file *os.File
	if file, err = os.Open(path); err != nil {
		return
	}

	defer file.Close()

	cfg = new(Config)
	var decoder = yaml.NewDecoder(file)
	if err = decoder.Decode(cfg); err != nil {
		return
	}

	return
}

func (x *Config) Validate() (err error) {
	if err = x.valPath(); err != nil {
		return
	}

	if err = x.valDriver(); err != nil {
		return
	}

	if err = x.valType(); err != nil {
		return
	}

	if err = x.valColumnType(); err != nil {
		return
	}

	if err = x.setModelNm(); err != nil {
		return
	}

	if err = x.valDatetime(); err != nil {
		return
	}

	if err = x.valIndex(); err != nil {
		return
	}

	if err = x.valRelation(); err != nil {
		return
	}

	return
}
