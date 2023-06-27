package parse

type (
	// ModelJson json format 정의
	ModelJson struct {
		TableNm string         `json:"tableNm"`
		Columns ModelColumnMap `json:"columns"`
		Indexes ModelIndexList `json:"indexes"`
	}
)

// ModelJson fields 구성요소
type (
	ModelColumn struct {
		Name          string `json:"name"`
		Gotype        string `json:"gotype"`
		Dbtype        string `json:"dbtype"`
		Unique        *bool  `json:"unique"`
		Primary       *bool  `json:"primary"`
		AutoIncrement *bool  `json:"autoIncrement"`
		Order         *Order `json:"order"`
	}

	ModelIndex struct {
		Columns []ColumnName `json:"columns"`
		Unique  *bool        `json:"unique"`
		Order   *Order       `json:"order"`
	}

	ModelIndexList []ModelIndex
	ModelColumnMap map[ColumnName]ModelColumn
	Order          string
	ColumnName     string
)

const (
	OrderASC  Order = "asc"
	OrderDESC Order = "desc"
)
