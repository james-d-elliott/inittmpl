package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/pelletier/go-toml"
	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "inittmpl <file>",

		Short: "Generate a configuration file from environment variables",
		Long:  `Generate a configuration file from environment variables.`,

		RunE: root,
		Args: cobra.ExactArgs(1),

		DisableAutoGenTag: true,
	}

	cmd.Flags().StringP("format", "f", "", "override the output format; yaml, toml, json")
	cmd.Flags().StringP("prefix", "p", "INITTMPL", "override the environment prefix")
	cmd.Flags().StringP("delimiter", "d", "__", "override the environment delimiter")
	cmd.Flags().BoolP("overwrite", "x", false, "overwrite existing file")
	cmd.Flags().BoolP("disable-lowercase", "e", false, "disable lowercase conversion of environment variables")
	cmd.Flags().BoolP("disable-opportunistic", "c", false, "disable opportunistic type conversion of environment variables")

	return cmd
}

func root(cmd *cobra.Command, args []string) (err error) {
	var (
		outputpath, prefix, delimiter, format string
		overwrite, lowercase, opportunistic   bool
	)

	if prefix, err = cmd.Flags().GetString("prefix"); err != nil {
		return err
	}

	if delimiter, err = cmd.Flags().GetString("delimiter"); err != nil {
		return err
	}

	if format, err = cmd.Flags().GetString("format"); err != nil {
		return err
	}

	if overwrite, err = cmd.Flags().GetBool("overwrite"); err != nil {
		return err
	}

	if lowercase, err = cmd.Flags().GetBool("disable-lowercase"); err != nil {
		return err
	}

	if opportunistic, err = cmd.Flags().GetBool("disable-opportunistic"); err != nil {
		return err
	}

	outputpath = args[0]

	if _, err = os.Stat(outputpath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if !overwrite {
		return nil
	}

	if format == "" {
		switch ext := filepath.Ext(outputpath); ext {
		case ".yaml", ".yml":
			format = "yaml"
		case ".toml":
			format = "toml"
		case ".json":
			format = "json"
		default:
			return fmt.Errorf("unknown output format for extension: %s", ext)
		}
	}

	k := koanf.New(".")

	cb := envCallback(!lowercase, !opportunistic, prefix, delimiter)

	if err = k.Load(env.ProviderWithValue(prefix, delimiter, cb), nil); err != nil {
		return fmt.Errorf("error occurred loading environment: %w", err)
	}

	values := map[string]any{}

	if err = k.Unmarshal("", &values); err != nil {
		return err
	}

	var f *os.File

	if f, err = os.OpenFile(outputpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return fmt.Errorf("error occurred during file open: %w", err)
	}

	defer f.Close()

	var encoder Encoder

	switch format {
	case "yaml":
		encoder = yaml.NewEncoder(f, yaml.Indent(2), yaml.UseSingleQuote(true), yaml.OmitEmpty())
	case "toml":
		e := toml.NewEncoder(f)

		e.Indentation("  ")

		encoder = e
	case "json":
		e := json.NewEncoder(f)

		e.SetIndent("", "  ")

		encoder = e
	default:
		return fmt.Errorf("error occurred during encoding: unknown output format: %s", format)
	}

	if err = encoder.Encode(values); err != nil {
		return fmt.Errorf("error occurred during encoding: %w", err)
	}

	return nil
}

func envCallback(lowercase, opportunistic bool, prefix, delimiter string) func(inkey, invalue string) (key string, value any) {
	return func(inkey, invalue string) (key string, value any) {
		var err error

		value, err = toValue(invalue, opportunistic)

		if err != nil {
			panic(fmt.Errorf("error occurred during explicit type conversion for env '%s': %w", inkey, err))
		}

		trimmed := strings.TrimPrefix(inkey, prefix+delimiter)

		if !lowercase {
			return trimmed, value
		}

		return strings.ToLower(trimmed), value
	}
}

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

type Encoder interface {
	Encode(v any) error
}
