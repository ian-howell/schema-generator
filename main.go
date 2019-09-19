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
	Use:           "schemagen file [file...]",
	Short:         "schemagen takes JSON/YAML and outputs a skeleton schema",
	Args:          cobra.MinimumNArgs(1),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		isDir, err := checkIfDir(o.outputLocation)
		if err != nil {
			return err
		}

		out := cmd.OutOrStdout()
		if isDir {
			err = os.MkdirAll(o.outputLocation, 0666)
			if err != nil {
				return err
			}
		} else if o.outputLocation != "" {
			// !isDir && o.outputLocation != ""
			out, err = os.Create(o.outputLocation)
			if err != nil {
				return err
			}
		}

		filenames := args
		for _, filename := range filenames {
			schemaName := generateSchemaName(filename)
			schema, err := generateSchemaForFile(filename, schemaName)
			if err != nil {
				return err
			}
			if isDir {
				err = outputToDirectory(schema, o.outputLocation, schemaName)
				if err != nil {
					return err
				}
			} else {
				_, err = out.Write(schema)
				if err != nil {
					return err
				}
			}
		}
		return nil
	},
}

type Options struct {
	outputLocation string
	indentLevel    int
}

var o Options

func init() {
	fs := rootCmd.PersistentFlags()
	fs.StringVarP(&o.outputLocation, "output", "o", "",
		"Output file. Providing a directory will output "+
			"one schema per input into the directory. "+
			"Leave empty to print to stdout")
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

func generateSchemaForFile(filename, schemaName string) ([]byte, error) {
	values, err := schemagen.ReadYAMLFile(filename)
	if err != nil {
		return []byte{}, err
	}
	schema := schemagen.GenerateSchema(schemaName, values)
	return schema.JSON(o.indentLevel)
}

func generateSchemaName(filename string) string {
	basename := filepath.Base(filename)
	return strings.Split(basename, ".")[0]
}

func checkIfDir(dirName string) (isDir bool, err error) {
	var fi os.FileInfo
	fi, err = os.Stat(dirName)
	isDir = (err == nil) && fi.IsDir()
	if os.IsNotExist(err) {
		// Non-existence is not an error
		err = nil
	}
	return isDir, err
}

func outputToDirectory(schema []byte, dirName, schemaName string) error {
	outputBasename := strings.Join([]string{schemaName, "schema", "json"}, ".")
	outputFilename := filepath.Join(dirName, outputBasename)
	return ioutil.WriteFile(outputFilename, schema, 0666)
}
