package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	schemagen "github.com/ian-howell/schema-generator/generator"
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(1)
	}

	filenames := os.Args[1:]
	for _, filename := range filenames {
		values, err := schemagen.ReadYAMLFile(filename)
		if err != nil {
			panic(err)
		}

		basename := filepath.Base(filename)
		name := strings.Split(basename, ".")[0]

		schema := schemagen.GenerateSchema(name, values)
		schemaJSON, err := schema.JSON(2)
		if err != nil {
			panic(err)
		}

		outputBasename := strings.Join([]string{name, "schema", "json"}, ".")
		outputFilename := filepath.Join(filepath.Dir(filename), outputBasename)
		ioutil.WriteFile(outputFilename, []byte(schemaJSON), 0644)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [filenames]\n", os.Args[0])
}
