package generator

import (
	"github.com/qdriven/qfluent-cli/pkg/inputs"
	"github.com/qdriven/qfluent-cli/pkg/log"
	"github.com/qdriven/qfluent-cli/pkg/transformer"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Generate is the main entry point for code generation/transformations.
func Generate(transformationsFile, source, destination string, inputArgs []string) error {
	transformations, err := transformer.Read(transformationsFile)
	if err != nil {
		return err
	}

	err = inputs.ParseCLIArgsInputs(transformations, inputArgs)
	if err != nil {
		return err
	}

	err = inputs.CollectUserInputs(transformations)
	if err != nil {
		return err
	}

	vars := collectSystemAndEnvironmentVariables(source, destination)
	err = transformations.Template(vars)
	if err != nil {
		return err
	}

	log.Debugf(spew.Sdump(transformations))
	err = transformer.Transform(source, destination, *transformations)
	if err != nil {
		return err
	}

	return nil
}

// Collects environment variables as well as system variables, e.g. source and destination
func collectSystemAndEnvironmentVariables(source, destination string) map[string]string {
	vars := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		vars[pair[0]] = vars[pair[1]]
	}
	vars["source"] = source
	vars["destination"] = destination
	return vars
}
