package scope

import (
	"github.com/orb-community/orb/pkg/errors"
	"github.com/orb-community/orb/pkg/types"
	"github.com/orb-community/orb/sinks/backend"
	"github.com/orb-community/orb/sinks/headers_type"
	"gopkg.in/yaml.v3"
)

const (
	HeadersConfigFeature = "X-Scope-OrgID"
)

var (
	features = []headers_type.ConfigFeature{
		{
			Type:     backend.ConfigFeatureTypeText,
			Input:    "text",
			Title:    "X-Scope-OrgID",
			Name:     HeadersConfigFeature,
			Required: true,
		},
	}
)

const HeaderType = "X-Scope-OrgID"

type HeaderConfig struct {
	Header string `json:"X-Scope-OrgID" ,yaml:"X-Scope-OrgID"`
}

func (a *HeaderConfig) Metadata() headers_type.HeadersTypeConfig {

	return headers_type.HeadersTypeConfig{
		Type:        HeaderType,
		Description: "Basic header",
		Config:      features,
	}
}

func (a *HeaderConfig) GetFeatureConfig() []headers_type.ConfigFeature {
	return features
}

func (a *HeaderConfig) ValidateConfiguration(inputFormat string, input interface{}) error {
	switch inputFormat {
	case "object":
		for key, value := range input.(types.Metadata) {
			vs := value.(string)
			if key == HeadersConfigFeature {
				if len(vs) == 0 {
					return errors.New("header cannot be empty")
				}
			}
		}
	case "yaml":
		err := yaml.Unmarshal([]byte(input.(string)), &a)
		if err != nil {
			return err
		}
		if len(a.Header) == 0 {
			return errors.New("header cannot be empty")
		}
	}
	return nil
}

func (a *HeaderConfig) ConfigToFormat(outputFormat string, input interface{}) (interface{}, error) {
	switch input.(type) {
	case types.Metadata:
		if outputFormat == "yaml" {
			retVal, err := yaml.Marshal(input)
			return string(retVal), err
		} else {
			return nil, errors.New("unsupported format")
		}
	case string:
		if outputFormat == "object" {
			retVal := make(types.Metadata)
			val := input.(string)
			err := yaml.Unmarshal([]byte(val), &retVal)
			return retVal, err
		} else {
			return nil, errors.New("unsupported format")
		}
	}
	return nil, errors.New("unsupported format")
}

func Register() {
	headerConfig := HeaderConfig{
		Header: "",
	}
	headers_type.Register(HeaderType, &headerConfig)
}
