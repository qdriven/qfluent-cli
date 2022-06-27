package transformer

import (
	"github.com/qdriven/qfluent-cli/pkg/inputs"
	"github.com/qdriven/qfluent-cli/pkg/operations"
)

const (
	TransformationTypeInclude = "include"
	TransformationTypeReplace = "replace"
)

type transformationsSpec struct {
	Ignore          []string             `yaml:"ignore"`
	Inputs          []inputs.InputSpec   `yaml:"inputs"`
	Transformations []transformationSpec `yaml:"transformations"`
	Before          operations.Spec      `yaml:"before"`
	After           operations.Spec      `yaml:"after"`
}

type transformationSpec struct {
	Name         string   `yaml:"name"`
	Type         string   `yaml:"type"`
	Pattern      string   `yaml:"pattern"`
	Replacement  string   `yaml:"replacement"`
	Files        []string `yaml:"files"`
	Condition    string   `yaml:"condition"`
	RegionMarker string   `yaml:"region_marker"`
}
