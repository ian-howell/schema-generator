package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	schemagen "github.com/ian-howell/schema-generator/generator"
)

var rootCmd = &cobra.Command{
	Use: "schemagen",
	Short: "schemagen takes JSON/YAML and outputs a skeleton schema",
	Args: cobra.MinimumNArgs(1),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filenames := args
		for _, filename := range filenames {
			values, err := schemagen.ReadYAMLFile(filename)
			if err != nil {
				return err
			}

			basename := filepath.Base(filename)
			name := strings.Split(basename, ".")[0]

			schema := schemagen.GenerateSchema(name, values)
			schemaJSON, err := schema.JSON(2)
			if err != nil {
				return err
			}

			outputBasename := strings.Join([]string{name, "schema", "json"}, ".")
			outputFilename := filepath.Join(filepath.Dir(filename), outputBasename)
			ioutil.WriteFile(outputFilename, []byte(schemaJSON), 0644)
		}
		return nil
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		out := rootCmd.OutOrStderr()
		fmt.Fprintf(out, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
