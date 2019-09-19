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
	Use:           "schemagen",
	Short:         "schemagen takes JSON/YAML and outputs a skeleton schema",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		filenames := args
		for _, filename := range filenames {
			if err := generateSchemaForFile(filename); err != nil {
				return err
			}
		}
		return nil
	},
}

type Options struct {
	indentLevel    int
}

var o Options

func init() {
	fs := rootCmd.PersistentFlags()
	fs.IntVarP(&o.indentLevel, "indent", "i", -1,
		"The number of spaces to indent at each level. "+
			"Default (-1) will print on a single line")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		out := rootCmd.OutOrStderr()
		fmt.Fprintf(out, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}

func generateSchemaForFile(filename string) error {
	values, err := schemagen.ReadYAMLFile(filename)
	if err != nil {
		return err
	}

	basename := filepath.Base(filename)
	name := strings.Split(basename, ".")[0]

	schema := schemagen.GenerateSchema(name, values)
	schemaJSON, err := schema.JSON(o.indentLevel)
	if err != nil {
		return err
	}

	outputBasename := strings.Join([]string{name, "schema", "json"}, ".")
	outputFilename := filepath.Join(filepath.Dir(filename), outputBasename)
	ioutil.WriteFile(outputFilename, []byte(schemaJSON), 0644)
	return nil
}
