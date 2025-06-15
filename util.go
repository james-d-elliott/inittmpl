package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/james-d-elliott/kissit/internal/consts"
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
	} else if strings.HasPrefix(in, "json::") {
		raw := strings.TrimPrefix(in, "json::")

		if strings.HasPrefix("[", raw) {
			v = []any{}
		} else {
			v = map[string]any{}
		}

		if err = json.Unmarshal([]byte(strings.TrimPrefix(in, "json::")), &v); err != nil {
			return nil, err
		}

		return v, nil
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
