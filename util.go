package main

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/pelletier/go-toml"

	"github.com/james-d-elliott/inittmpl/internal/consts"
)

func toValue(in string, opportunistic bool) (v any, err error) {
	if strings.HasPrefix(in, "string::") {
		return strings.TrimPrefix(in, "string::"), nil
	} else if strings.HasPrefix(in, "int::") {
		return strconv.ParseInt(strings.TrimPrefix(in, "int::"), 10, 64)
	} else if strings.HasPrefix(in, "uint::") {
		return strconv.ParseUint(strings.TrimPrefix(in, "uint::"), 10, 64)
	} else if strings.HasPrefix(in, "bool::") {
		return strconv.ParseBool(strings.TrimPrefix(in, "bool::"))
	} else if strings.HasPrefix(in, "float::") {
		return strconv.ParseFloat(strings.TrimPrefix(in, "float::"), 10)
	}

	if opportunistic {
		if v, err = strconv.Atoi(in); err == nil {
			return v, nil
		} else if v, err = strconv.ParseBool(in); err == nil {
			return v, nil
		}
	}

	return in, nil
}

func getEncoder(format string) (f func(wr io.Writer) Encoder, err error) {
	switch format {
	case consts.FormatYAML:
		return func(wr io.Writer) Encoder {
			return yaml.NewEncoder(wr, yaml.Indent(2), yaml.UseSingleQuote(true), yaml.OmitEmpty())
		}, nil
	case consts.FormatTOML:
		return func(wr io.Writer) Encoder {
			e := toml.NewEncoder(wr)
			e.Indentation("  ")

			return e
		}, nil
	case consts.FormatJSON:
		return func(wr io.Writer) Encoder {
			encoder := json.NewEncoder(wr)

			encoder.SetIndent("", "  ")

			return encoder
		}, nil
	default:
		return nil, fmt.Errorf("error occurred during encoding: unknown output format: %s", format)
	}
}

func extToFormat(ext, fallback string) string {
	switch ext {
	case ".yaml", ".yml":
		return consts.FormatYAML
	case ".toml":
		return consts.FormatTOML
	case ".json":
		return consts.FormatJSON
	default:
		return fallback
	}
}
