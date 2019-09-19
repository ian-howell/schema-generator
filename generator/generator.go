package generator

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"sigs.k8s.io/yaml"
)

// Schema represents the document structure to validate the values.yaml file against
type Schema map[string]interface{}

// JSON encodes the Schema into JSON.
func (s Schema) JSON(indent int) ([]byte, error) {
	if indent >= 0 {
		return json.MarshalIndent(s, "", strings.Repeat(" ", indent))
	}
	return json.Marshal(s)
}

// ReadYAML will parse YAML byte data into a map[string]interface{}.
func ReadYAML(data []byte) (vals map[string]interface{}, err error) {
	err = yaml.Unmarshal(data, &vals)
	if len(vals) == 0 {
		vals = map[string]interface{}{}
	}
	return vals, err
}

// ReadYAMLFile will parse a YAML file into a map[string]interface{}.
func ReadYAMLFile(filename string) (map[string]interface{}, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return ReadYAML(data)
}

// GenerateSchema will create a JSON Schema (in YAML format) for the given values
func GenerateSchema(name string, values map[string]interface{}) Schema {
	schema := Schema{
		gojsonschema.KEY_TYPE:  gojsonschema.TYPE_OBJECT,
		gojsonschema.KEY_TITLE: name,
	}
	if len(values) > 0 {
		schema[gojsonschema.STRING_PROPERTIES] = parsePropertiesFromYAML(values)
	}
	return schema
}

func parsePropertiesFromYAML(values map[string]interface{}) map[string]map[string]interface{} {
	properties := make(map[string]map[string]interface{})
	for k, v := range values {
		// If the value is null, then there's no way to determine the properties attributes
		if v == nil || v == "" {
			continue
		}

		properties[k] = make(map[string]interface{})
		// the following types are the only types possible from unmarshalling
		switch v := v.(type) {
		case bool:
			properties[k][gojsonschema.KEY_TYPE] = gojsonschema.TYPE_BOOLEAN
		case float64:
			properties[k][gojsonschema.KEY_TYPE] = gojsonschema.TYPE_NUMBER
		case string:
			properties[k][gojsonschema.KEY_TYPE] = gojsonschema.TYPE_STRING
		case []interface{}:
			properties[k][gojsonschema.KEY_TYPE] = gojsonschema.TYPE_ARRAY
			properties[k][gojsonschema.KEY_ITEMS] = parseItemsFromYAML(v)
		case map[string]interface{}:
			properties[k][gojsonschema.KEY_TYPE] = gojsonschema.TYPE_OBJECT
			object := parsePropertiesFromYAML(v)
			if len(object) > 0 {
				properties[k][gojsonschema.TYPE_OBJECT] = object
			}
		}
	}
	return properties
}

func parseItemsFromYAML(items []interface{}) map[string]interface{} {
	properties := make(map[string]interface{})
	if len(items) == 0 {
		return properties
	}
	v := items[0]
	// the following types are the only types possible from unmarshalling
	switch v := v.(type) {
	case bool:
		properties[gojsonschema.KEY_TYPE] = gojsonschema.TYPE_BOOLEAN
	case float64:
		properties[gojsonschema.KEY_TYPE] = gojsonschema.TYPE_NUMBER
	case string:
		properties[gojsonschema.KEY_TYPE] = gojsonschema.TYPE_STRING
	case []interface{}:
		properties[gojsonschema.KEY_TYPE] = gojsonschema.TYPE_ARRAY
		properties[gojsonschema.KEY_ITEMS] = parseItemsFromYAML(v)
	case map[string]interface{}:
		properties[gojsonschema.KEY_TYPE] = gojsonschema.TYPE_OBJECT
		object := parsePropertiesFromYAML(v)
		if len(object) > 0 {
			properties[gojsonschema.TYPE_OBJECT] = object
		}
	}
	return properties
}
