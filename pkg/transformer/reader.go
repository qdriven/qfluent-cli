package transformer

import (
	"github.com/qdriven/qfluent-cli/pkg/inputs"
	types "github.com/qdriven/qfluent-cli/pkg/io"
	"github.com/qdriven/qfluent-cli/pkg/operations"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func Read(transformationsFile string) (*Transformations, error) {
	yamlFile, err := ioutil.ReadFile(transformationsFile)
	if err != nil {
		return nil, err
	}
	var spec transformationsSpec
	err = yaml.Unmarshal(yamlFile, &spec)
	if err != nil {
		return nil, err
	}
	return FromSpec(spec)
}

func FromSpec(spec transformationsSpec) (*Transformations, error) {
	return &Transformations{
		ignore:       types.NewFilePatterns(spec.Ignore),
		transformers: transformersFromSpec(spec.Transformations),
		prompters:    inputs.FromSpec(spec.Inputs),
		userInputs:   make(map[string]inputs.PromptResponse),
		before:       operations.FromSpec(spec.Before),
		after:        operations.FromSpec(spec.After),
	}, nil
}

func transformersFromSpec(transformationSpecs []transformationSpec) []Transformer {
	var transformers []Transformer
	for _, t := range transformationSpecs {
		transformers = append(transformers, newTransformer(t))
	}
	return transformers
}
