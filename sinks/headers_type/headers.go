package headers_type

type HeadersType interface {
	Metadata() HeadersTypeConfig
	GetFeatureConfig() []ConfigFeature
	ValidateConfiguration(inputFormat string, input interface{}) error
	ConfigToFormat(outputFormat string, input interface{}) (interface{}, error)


}

const HeadersKey = "headers"

type ConfigFeature struct {
	Type     string `json:"type"`
	Input    string `json:"input"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type HeadersTypeConfig struct {
	Type        string          `json:"type"`
	Description string          `json:"description"`
	Config      []ConfigFeature `json:"config"`
}

var headersTypes = make(map[string]HeadersType)

func Register(name string, headersType HeadersType) {
	headersTypes[name] = headersType
}

func GetList() []HeadersTypeConfig {
	keys := make([]HeadersTypeConfig, 0, len(headersTypes))
	for _, v := range headersTypes {
		keys = append(keys, v.Metadata())
	}
	return keys
}

func GetHeadersType(id string) (HeadersType, bool) {
	v, ok := headersTypes[id]
	return v, ok
}
