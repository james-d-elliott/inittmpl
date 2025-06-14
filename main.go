package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
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

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
