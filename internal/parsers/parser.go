package parsers

import (
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/v2"

	"github.com/james-d-elliott/inittmpl/internal/consts"
	"github.com/james-d-elliott/inittmpl/internal/parsers/yaml"
)

func Parser(format string) koanf.Parser {
	switch format {
	case consts.FormatYAML:
		return yaml.Parser()
	case consts.FormatTOML:
		return toml.Parser()
	case consts.FormatJSON:
		return json.Parser()
	default:
		return nil
	}
}
