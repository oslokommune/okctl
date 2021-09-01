package store

import (
	"encoding/json"
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"

	"sigs.k8s.io/yaml"
)

type toJSON struct{}

func (t *toJSON) PreProcess(data interface{}) (*PreProcessed, error) {
	d, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, fmt.Errorf(constant.MarshalJsonError, err)
	}

	return &PreProcessed{
		Format: t.Type(),
		Data:   d,
	}, nil
}

func (t *toJSON) Type() string {
	return "json"
}

// ToJSON provides a preprocessor that marshals a struct to json format
func ToJSON() PreProcessor {
	return &toJSON{}
}

type fromJSON struct{}

func (f *fromJSON) PostProcess(into interface{}, data []byte) (*PostProcessed, error) {
	err := json.Unmarshal(data, into)
	if err != nil {
		return nil, fmt.Errorf(constant.MarshalJsonError, err)
	}

	return &PostProcessed{
		Format: f.Type(),
		Data:   into,
	}, nil
}

func (f *fromJSON) Type() string {
	return "json"
}

// FromJSON provides a postprocessor that unmarshals from yaml into a struct
func FromJSON() PostProcessor {
	return &fromJSON{}
}

type toYAML struct{}

func (t *toYAML) PreProcess(data interface{}) (*PreProcessed, error) {
	d, err := yaml.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf(constant.MarshalYamlError, err)
	}

	return &PreProcessed{
		Format: t.Type(),
		Data:   d,
	}, nil
}

func (t *toYAML) Type() string {
	return "yaml"
}

// ToYAML provides a preprocessor that marshals a struct to yaml format
func ToYAML() PreProcessor {
	return &toYAML{}
}

type fromYAML struct{}

func (f *fromYAML) PostProcess(into interface{}, data []byte) (*PostProcessed, error) {
	err := yaml.Unmarshal(data, into)
	if err != nil {
		return nil, fmt.Errorf(constant.MarshalYamlError, err)
	}

	return &PostProcessed{
		Format: f.Type(),
		Data:   into,
	}, nil
}

func (f *fromYAML) Type() string {
	return "yaml"
}

// FromYAML provides a postprocessor that unmarshals yaml data ino a struct
func FromYAML() PostProcessor {
	return &fromYAML{}
}
