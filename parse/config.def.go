package parse

type (
	// ConfigJson 접속정보, 사용할 모델 정보
	ConfigJson struct {
		Connection *FConnection `json:"connection" yaml:"connection"`
		Models     FModelMap    `json:"models"`
	}
)

// ConfigJson 구성요소
type (
	FConnection struct {
		Host     string
		Username string
		Password string
	}

	FModel struct {
		Binds FModelBindList `json:"binds" yaml:"binds"`
	}

	FModelBind struct {
		As            string   `json:"as" yaml:"as"`
		From          string   `json:"from" yaml:"from"`
		ForeignColumn string   `json:"foreignColumn" yaml:"foreignColumn"`
		LocalColumn   string   `json:"localColumn" yaml:"localColumn"`
		BindType      BindType `json:"bindType" yaml:"bindType"`
	}

	FModelMap      map[ModelName]*FModel
	FModelBindList []*FModelBind
	ModelName      string
	BindType       string
)

const (
	BindTypeOne  BindType = "one"
	BindTypeList BindType = "list"
)
