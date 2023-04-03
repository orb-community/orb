package authentication_type

type AuthenticationType interface {
	Metadata() interface{}
	GetFeatureConfig() []ConfigFeature
	ValidateConfiguration(inputFormat string, input interface{}) error
	ConfigToFormat(outputFormat string, input interface{}) (interface{}, error)
	OmitInformation(outputFormat string, input interface{}) (interface{}, error)
	EncodeInformation(outputFormat string, input interface{}) (interface{}, error)
	DecodeInformation(outputFormat string, input interface{}) (interface{}, error)
}

type ConfigFeature struct {
	Type     string `json:"type"`
	Input    string `json:"input"`
	Title    string `json:"title"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

var authTypes = make(map[string]AuthenticationType)

func Register(name string, b AuthenticationType) {
	authTypes[name] = b
}

func GetList() []string {
	keys := make([]string, 0, len(authTypes))
	for k := range authTypes {
		keys = append(keys, k)
	}
	return keys
}

func GetAuthType(id string) (AuthenticationType, bool) {
	v, ok := authTypes[id]
	return v, ok
}
