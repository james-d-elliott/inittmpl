package yaml

import (
	"bytes"
	"github.com/goccy/go-yaml"
)

// YAML implements a YAML parser.
type YAML struct{}

// Parser returns a YAML Parser.
func Parser() *YAML {
	return &YAML{}
}

// Unmarshal parses the given YAML bytes.
func (p *YAML) Unmarshal(b []byte) (map[string]any, error) {
	var out map[string]any
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}

	return out, nil
}

// Marshal marshals the given config map to YAML bytes.
func (p *YAML) Marshal(o map[string]any) (data []byte, err error) {
	buf := &bytes.Buffer{}

	if err = yaml.NewEncoder(buf, yaml.Indent(2), yaml.UseSingleQuote(true), yaml.OmitEmpty()).Encode(o); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
