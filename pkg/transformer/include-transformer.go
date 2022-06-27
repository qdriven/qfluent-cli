package transformer

import (
	"bufio"
	"fmt"
	ioutils "github.com/qdriven/qfluent-cli/pkg/ioutils"
	qlog "github.com/qdriven/qfluent-cli/pkg/log"
	"github.com/qdriven/qfluent-cli/pkg/template"
	"strings"
)

type includeTransformer struct {
	name string
	// the condition as a string
	condition string
	// regions marker in the file
	regionMarker string
	files        []ioutils.FilePattern
	// Does this condition evaluate to true, provided the variable values?
	truthy bool
}

func newIncludeTransformer(spec transformationSpec) *includeTransformer {
	return &includeTransformer{
		name:         spec.Name,
		condition:    spec.Condition,
		regionMarker: spec.RegionMarker,
		files:        ioutils.NewFilePatterns(spec.Files),
	}
}

func (t *includeTransformer) GetName() string {
	return t.name
}

func (t *includeTransformer) GetFilePatterns() []ioutils.FilePattern {
	return t.files
}

func (t *includeTransformer) Transform(input ioutils.File) ioutils.File {
	if len(t.regionMarker) == 0 {
		if t.truthy {
			return input
		}
		// Discard the entire file
		return ioutils.File{
			Discarded:    true,
			FullPath:     input.FullPath,
			RelativePath: input.RelativePath,
		}
	}
	// Locate begin and end lines of the markers
	beginMarker := fmt.Sprintf("BEGIN %s", t.regionMarker)
	endMarker := fmt.Sprintf("END %s", t.regionMarker)
	scanner := bufio.NewScanner(strings.NewReader(input.Contents))
	var (
		output      strings.Builder
		insideBlock = false
	)
	for scanner.Scan() {
		includeLine := t.truthy || !insideBlock
		text := scanner.Text()
		if strings.Contains(text, beginMarker) {
			insideBlock = true
			includeLine = false
		}
		if strings.Contains(text, endMarker) {
			insideBlock = false
			includeLine = false
		}
		if includeLine {
			output.WriteString(text)
			output.WriteRune('\n')
		}
	}
	if scanner.Err() != nil {
		qlog.Errorf("Error while scanning file %s: %+v.\n\n Contents: %s ...",
			scanner.Err(), input.FullPath, input.Contents[:100])
	}

	newContents := output.String()
	// Check if a the last newline should be preserved or discarded.
	if len(newContents) > 0 && !t.hasEmptyLineAtTheEnd(input.Contents) {
		newContents = newContents[:len(newContents)-1]
	}

	return ioutils.File{
		Contents:     newContents,
		FullPath:     input.FullPath,
		RelativePath: input.RelativePath,
	}
}

func (t *includeTransformer) Template(vars map[string]string) error {
	var err error
	t.truthy, err = template.EvaluateCondition(t.condition, vars)
	return err
}

func (t *includeTransformer) hasEmptyLineAtTheEnd(s string) bool {
	l := len(s)
	return l >= 1 && s[l-1] == '\n'
}

var _ Transformer = &includeTransformer{}
