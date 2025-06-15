package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/cobra"

	"github.com/james-d-elliott/inittmpl/internal/parsers"
)

func newRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "inittmpl <file>",

		Short:   "Generate a configuration file from environment variables",
		Long:    rootLong,
		Example: rootExample,
		RunE:    root,
		Args:    cobra.ExactArgs(1),

		DisableAutoGenTag: true,
	}

	cmd.Flags().StringP("format", "f", "", "override the output format; yaml, toml, json")
	cmd.Flags().StringP("prefix", "p", "INITTMPL", "override the environment prefix")
	cmd.Flags().StringP("delimiter", "d", "__", "override the environment delimiter")
	cmd.Flags().BoolP("overwrite", "x", false, "overwrite existing file")
	cmd.Flags().StringSliceP("files", "z", nil, "uses the specified files as additional input, if any of these are the same file as the output assumes overwrite")
	cmd.Flags().BoolP("disable-lowercase", "e", false, "disable lowercase conversion of environment variables")
	cmd.Flags().BoolP("disable-opportunistic", "c", false, "disable opportunistic type conversion of environment variables")

	return cmd
}

func root(cmd *cobra.Command, args []string) (err error) {
	var (
		outputpath                          string
		inputpaths                          []string
		prefix, delimiter, format           string
		overwrite, lowercase, opportunistic bool
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

	input := false

	if input = cmd.Flags().Changed("files"); input {
		if inputpaths, err = cmd.Flags().GetStringSlice("files"); err != nil {
			return err
		}

		outputabs, _ := filepath.Abs(outputpath)

		if !overwrite {
			var inputabs string

			for _, inputpath := range inputpaths {

				if inputabs, err = filepath.Abs(inputpath); err == nil && inputabs == outputabs {
					overwrite = true

					break
				}
			}
		}
	}

	if _, err = os.Stat(outputpath); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else if !overwrite {
		return nil
	}

	if format == "" {
		ext := filepath.Ext(outputpath)

		format = extToFormat(ext, format)

		if format == "" {
			return fmt.Errorf("unknown output format for extension: %s", ext)
		}
	}

	k := koanf.New(".")

	if input {
		for _, inputpath := range inputpaths {
			ext := filepath.Ext(inputpath)

			informat := extToFormat(ext, format)

			if err = k.Load(file.Provider(inputpath), parsers.Parser(informat)); err != nil {
				return fmt.Errorf("error loading input file '%s': %w", inputpath, err)
			}
		}
	}

	cb := envCallback(!lowercase, !opportunistic, prefix, delimiter)

	if err = k.Load(env.ProviderWithValue(prefix, delimiter, cb), nil); err != nil {
		return fmt.Errorf("error occurred loading environment: %w", err)
	}

	var f *os.File

	if f, err = os.OpenFile(outputpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); err != nil {
		return fmt.Errorf("error occurred during file open: %w", err)
	}

	defer f.Close()

	var data []byte

	if data, err = k.Marshal(parsers.Parser(format)); err != nil {
		return fmt.Errorf("error occurred during file marshal: %w", err)
	}

	if _, err = f.Write(data); err != nil {
		return fmt.Errorf("error occurred during file write: %w", err)
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

type Encoder interface {
	Encode(v any) error
}
