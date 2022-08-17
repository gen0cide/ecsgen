package gostruct

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/gen0cide/ecsgen"
	"github.com/gen0cide/ecsgen/generator"
	"github.com/urfave/cli/v2"
	"golang.org/x/tools/imports"
)

var (
	// ErrInvalidPackageName is thrown when a Go package name is either not specified or is not valid.
	ErrInvalidPackageName = errors.New("package name was either empty or an invalid go package identifier")

	// ErrInvalidOutputDir is thrown when the output directory does not exist.
	ErrInvalidOutputDir = errors.New("output directory was either blank or did not exist")
)

var defaultFilename = "generated_ecs.go"

type basic struct {
	PackageName        string
	OutputDir          string
	Filename           string
	IncludeJSONMarshal bool
	WeakUnmarshal      bool
}

// New is a constructor for an empty debug output plugin.
func New() generator.Generator {
	return &basic{}
}

// ID implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (b *basic) ID() string {
	return "gostruct"
}

// CLIFlags implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (b *basic) CLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "package-name",
			Usage:       "Name of the Go package for the generated code.",
			EnvVars:     []string{"PACKAGE_NAME"},
			Destination: &b.PackageName,
		},
		&cli.StringFlag{
			Name:        "output-dir",
			Usage:       "Path to the directory where the generated code should be written.",
			EnvVars:     []string{"OUTPUT_DIR"},
			Destination: &b.OutputDir,
		},
		&cli.StringFlag{
			Name:        "output-filename",
			Usage:       fmt.Sprintf("Destination filename for the generated code. (default: %s)", defaultFilename),
			EnvVars:     []string{"OUTPUT_FILENAME"},
			Destination: &b.Filename,
		},
		&cli.BoolFlag{
			Name:        "marshal-json",
			Usage:       "Include a json.Marshaler implementation that removes empty fields.",
			EnvVars:     []string{"MARSHAL_JSON"},
			Destination: &b.IncludeJSONMarshal,
		},
		&cli.BoolFlag{
			Name:        "weak-unmarshal",
			Usage:       "If set, when unmarshalling objects it will convert scalar types to arrays, non numeric types to numeric and vice versa.",
			EnvVars:     []string{"WEAK_UNMARSHAL"},
			Destination: &b.WeakUnmarshal,
		},
	}
}

// Validate implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (b *basic) Validate() error {
	// Check the Output Directory
	// is it assigned?
	if b.OutputDir == "" {
		return ErrInvalidOutputDir
	}

	// Use the default unless otherwise specified
	if b.Filename == "" {
		b.Filename = defaultFilename
	}

	// Is it a valid path?
	dir, err := os.Stat(b.OutputDir)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrInvalidOutputDir
		}

		return fmt.Errorf("error locating specified output directory: %v", err)
	}

	// is it a valid directory?
	if !dir.IsDir() {
		return fmt.Errorf("specified output directory was a path to a file, not a directory")
	}

	// while Go maintains STRONG guidance on package naming conventions,
	// it doesn't actually seem to enforce a whole lot. Keeping it basic for now.
	pkgRegex := regexp.MustCompile(`^[a-zA-Z0-9\_]{1,64}$`)
	if !pkgRegex.MatchString(b.PackageName) {
		return ErrInvalidPackageName
	}

	return nil
}

// GoFieldType returns the Go type to be used in the Go struct field type definition.
func GoFieldType(n *ecsgen.Node) string {
	// create a buffer to determine type
	typeBuf := new(bytes.Buffer)

	// add array syntax if the field normalizes out to an array
	if n.IsArray() {
		typeBuf.WriteString("[]")
	}

	// if Node is an Object, we need to return this object's type. For example,
	// Node("client.nat") needs to return "ClientNAT" as it's Go type.
	if n.IsObject() {
		typeBuf.WriteString(n.TypeIdent().Pascal())
		return typeBuf.String()
	}

	// Special cases denoted by the ECS developers.
	switch {
	case n.Name == "duration" && n.Definition.Type == "long":
		typeBuf.WriteString("time.Duration")
		return typeBuf.String()
	case n.Path == "labels":
		typeBuf.WriteString("map[string]interface{}")
		return typeBuf.String()
	}

	// Find the right type!
	switch n.Definition.Type {
	case "keyword", "text", "ip", "geo_point", "constant_keyword", "match_only_text", "wildcard":
		typeBuf.WriteString("string")
		return typeBuf.String()
	case "long":
		typeBuf.WriteString("int64")
		return typeBuf.String()
	case "integer":
		typeBuf.WriteString("int32")
		return typeBuf.String()
	case "float", "scaled_float":
		typeBuf.WriteString("float64")
		return typeBuf.String()
	case "date":
		typeBuf.WriteString("time.Time")
		return typeBuf.String()
	case "boolean":
		typeBuf.WriteString("bool")
		return typeBuf.String()
	case "object", "nested":
		typeBuf.WriteString("map[string]interface{}")
		return typeBuf.String()
	case "flattened":
		typeBuf.WriteString("map[string]string")
		return typeBuf.String()
	default:
		panic(fmt.Errorf("no translation for %v (field %s)", n.Definition.Type, n.Name))
	}
}

// ToGoCode attempts to convert an ecsgen.Node into a Golang struct definition.
func (b *basic) ToGoCode(n *ecsgen.Node) (string, error) {
	// we can only generate a Go struct definition for an Object, verify
	// we're not shooting ourselves in the foot
	if !n.IsObject() {
		return "", fmt.Errorf("node %s is not an object", n.Path)
	}

	// Now enumerate the Node's fields and sort the keys so the resulting Go code
	// is deterministically generated
	fieldKeys := []string{}

	for key := range n.Children {
		fieldKeys = append(fieldKeys, key)
	}

	sort.Strings(fieldKeys)

	// Create a new buffer to write the struct definition to
	buf := new(strings.Builder)

	// comment and type definition
	buf.WriteString(fmt.Sprintf("// %s defines the object located at ECS path %s.", n.TypeIdent().Pascal(), n.Path))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("type %s struct {", n.TypeIdent().Pascal()))
	buf.WriteString("\n")

	// Enumerate the fields and generate their field definition, adding it
	// to the buffer as a line item.
	for _, k := range fieldKeys {
		scalarField := n.Children[k]
		buf.WriteString(
			fmt.Sprintf(
				"\t%s %s `json:\"%s,omitempty\" yaml:\"%s,omitempty\" ecs:\"%s\"`",
				scalarField.FieldIdent().Pascal(),
				GoFieldType(scalarField),
				scalarField.Name,
				scalarField.Name,
				scalarField.Path,
			),
		)
		buf.WriteString("\n")
	}

	// Close the type definition and return the result
	buf.WriteString("}")
	buf.WriteString("\n")

	// if the user included the JSON operator flag, add the implementation
	if b.IncludeJSONMarshal {
		// Now we implement at json.Marshaler implementation for each specific type that
		// removes any nested JSON types that might exist.
		//
		// We do this by enumerating every field in the type and check to see
		// if it's got a zero value.
		buf.WriteString("\n")
		buf.WriteString("// MarshalJSON implements the json.Marshaler interface and removes zero values from returned JSON.")
		buf.WriteString("\n")
		buf.WriteString(
			fmt.Sprintf(
				"func (b %s) MarshalJSON() ([]byte, error) {",
				n.TypeIdent().Pascal(),
			),
		)
		buf.WriteString("\n")

		// Define the result struct we will populate non-zero fields with
		buf.WriteString("\tres := map[string]interface{}{}")
		buf.WriteString("\n")
		buf.WriteString("\n")

		// enumerate the fields for the object fields
		for _, fieldName := range fieldKeys {
			field := n.Children[fieldName]
			if GoFieldType(field) != "bool" {
				buf.WriteString(
					fmt.Sprintf(
						"\tif val := reflect.ValueOf(b.%s); !val.IsZero() {", field.FieldIdent().Pascal(),
					),
				)
			}
			buf.WriteString(
				fmt.Sprintf(
					"\t\tres[\"%s\"] = b.%s",
					field.Name,
					field.FieldIdent().Pascal(),
				),
			)
			if GoFieldType(field) != "bool" {
				buf.WriteString("\t}")
			}
			buf.WriteString("\n")
			buf.WriteString("\n")
		}

		// add a line spacer and return the marshaled JSON result
		buf.WriteString("\n")
		buf.WriteString("\treturn json.Marshal(res)")
		buf.WriteString("\n")
		buf.WriteString("}")
		buf.WriteString("\n")
	}

	if b.WeakUnmarshal {
		// We override the UnmarshalJSON method to be able to attempt to unmarshal scalar fields that should be
		// an array and to convert strings that should be put into numeric fields
		buf.WriteString("\n")
		buf.WriteString("// UnmarshalJSON implements the json.Unmarshaler interface and attempts weak decoding of fields")
		buf.WriteString("\n")
		buf.WriteString(
			fmt.Sprintf(
				"func (b *%s) UnmarshalJSON(data []byte) error {",
				n.TypeIdent().Pascal(),
			),
		)
		buf.WriteString("\n")

		for i, fieldName := range fieldKeys {
			field := n.Children[fieldName]
			if i == 0 {
				buf.WriteString(
					fmt.Sprintf(
						"\tresult := gjson.GetBytes(data, \"%s\")\n", field.Name,
					),
				)
			} else {
				buf.WriteString(
					fmt.Sprintf(
						"\tresult = gjson.GetBytes(data, \"%s\")\n", field.Name,
					),
				)
			}
			buf.WriteString("\tif result.Index > 0 {\n")
			if GoFieldType(field) == "[]string" || GoFieldType(field) == "int32" || GoFieldType(field) == "int64" {
				if GoFieldType(field) == "[]string" {
					buf.WriteString(
						fmt.Sprintf(
							"\t\tval_%s, err := decodeArray(\"%s\", result.Value())\n", field.Name, field.Path,
						),
					)
				} else if GoFieldType(field) == "int32" {
					buf.WriteString(
						fmt.Sprintf(
							"\t\tval_%s, err := decodeInt32(\"%s\", result.Value())\n", field.Name, field.Path,
						),
					)
				} else if GoFieldType(field) == "int64" {
					buf.WriteString(
						fmt.Sprintf(
							"\t\tval_%s, err := decodeInt64(\"%s\", result.Value())\n", field.Name, field.Path,
						),
					)
				}
				buf.WriteString("\t\tif err != nil {\n")
				buf.WriteString("\t\t\treturn err\n")
				buf.WriteString("\t\t}\n")
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = val_%s\n", field.FieldIdent().Pascal(), field.Name,
					),
				)
			} else if GoFieldType(field) == "string" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.String()\n", field.FieldIdent().Pascal(),
					),
				)
			} else if GoFieldType(field) == "time.Time" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.Time()\n", field.FieldIdent().Pascal(),
					),
				)
			} else if GoFieldType(field) == "float64" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.Float()\n", field.FieldIdent().Pascal(),
					),
				)
			} else if GoFieldType(field) == "bool" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.Bool()\n", field.FieldIdent().Pascal(),
					),
				)
			} else {
				buf.WriteString("\t\tvar raw []byte\n")
				buf.WriteString("\t\traw = data[result.Index:result.Index+len(result.Raw)]\n")
				buf.WriteString(fmt.Sprintf("\t\tvar aux %s\n", GoFieldType(field)))
				buf.WriteString("\t\terr := json.Unmarshal(raw, &aux)\n")
				buf.WriteString("\t\tif err != nil {\n")
				buf.WriteString("\t\t\treturn err\n")
				buf.WriteString("\t\t}\n")
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = aux\n", field.FieldIdent().Pascal(),
					),
				)
			}
			buf.WriteString("\t}\n")
			buf.WriteString("\n")
		}

		buf.WriteString("\treturn nil\n")
		buf.WriteString("}")
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// CreateBase generates the top level ECS Base struct that holds all fieldsets and top level fields.
func (b *basic) CreateBase(r *ecsgen.Root) (string, error) {
	// buckets to sort the field names into
	scalarFields := []string{}
	objectFields := []string{}

	// first we need to sort the field names, and separate out Base fields
	// from the FieldSets
	for fieldName, fieldNode := range r.TopLevel {
		if fieldNode.IsObject() {
			objectFields = append(objectFields, fieldName)
			continue
		}

		scalarFields = append(scalarFields, fieldName)
	}

	sort.Strings(scalarFields)
	sort.Strings(objectFields)

	// now to build the buffer that holds the Go type definition
	buf := new(strings.Builder)

	// Add the type comment and the definition to the buffer
	buf.WriteString("// Base defines the top level Elastic Common Schema (ECS) type. This type should be the default for interacting with ECS data, including the marshaling and unmarshaling of it.")
	buf.WriteString("\n")
	buf.WriteString("type Base struct {")
	buf.WriteString("\n")

	// Enumerate the scalar fields (the fields that are direct types in the Base fieldset)
	// and add them to the type definition
	for _, k := range scalarFields {
		field := r.TopLevel[k]
		buf.WriteString(
			fmt.Sprintf(
				"\t%s %s `json:\"%s,omitempty\" yaml:\"%s,omitempty\" ecs:\"%s\"`",
				field.FieldIdent().Pascal(),
				GoFieldType(field),

				// We don't actually use the "parsed field name" here because
				// unfortunately we have to account for the @timestamp field name
				// because YOLO, that field follows other naming conventions!
				k,
				k,
				k,
			),
		)
		buf.WriteString("\n")
	}

	// Now enumerate the object fields and add those to the base type
	for _, k := range objectFields {
		field := r.TopLevel[k]
		buf.WriteString(
			fmt.Sprintf(
				"\t%s %s `json:\"%s,omitempty\" yaml:\"%s,omitempty\" ecs:\"%s\"`",
				field.FieldIdent().Pascal(),
				GoFieldType(field),
				field.Name,
				field.Name,
				field.Path,
			),
		)
		buf.WriteString("\n")
	}

	// close the struct
	buf.WriteString("}")
	buf.WriteString("\n")

	// if the user indicated they wanted a json.Marshaler implementation,
	// then generate that.
	if b.IncludeJSONMarshal {
		// Now we have to create the marshaler to account for Zero values!
		// this will remove object fields that are empty from the resulting JSON.
		//
		// The way we do this is by enumerating every field in the top level Base
		// and check to see if it's got a zero value.
		buf.WriteString("\n")
		buf.WriteString("// MarshalJSON implements the json.Marshaler interface and removes zero values from returned JSON.")
		buf.WriteString("\n")
		buf.WriteString("func (b Base) MarshalJSON() ([]byte, error) {")
		buf.WriteString("\n")

		// Define the result struct we will populate non-zero fields with
		buf.WriteString("\tres := map[string]interface{}{}")
		buf.WriteString("\n")
		buf.WriteString("\n")

		// first we enumerate the scalar fields
		for _, fieldName := range scalarFields {
			field := r.TopLevel[fieldName]
			if GoFieldType(field) != "bool" {
				buf.WriteString(
					fmt.Sprintf(
						"\tif val := reflect.ValueOf(b.%s); !val.IsZero() {", field.FieldIdent().Pascal(),
					),
				)
			}
			buf.WriteString(
				fmt.Sprintf(
					"\t\tres[\"%s\"] = b.%s",
					fieldName,
					field.FieldIdent().Pascal(),
				),
			)
			if GoFieldType(field) != "bool" {
				buf.WriteString("\t}")
			}
			buf.WriteString("\n")
			buf.WriteString("\n")
		}

		// now we enumerate the object fields
		for _, fieldName := range objectFields {
			field := r.TopLevel[fieldName]
			buf.WriteString(
				fmt.Sprintf(
					"\tif val := reflect.ValueOf(b.%s); !val.IsZero() {", field.FieldIdent().Pascal(),
				),
			)
			buf.WriteString(
				fmt.Sprintf(
					"\t\tres[\"%s\"] = b.%s",
					field.Name,
					field.FieldIdent().Pascal(),
				),
			)
			buf.WriteString("\t}")
			buf.WriteString("\n")
			buf.WriteString("\n")
		}

		// add a line spacer and return the marshaled JSON result
		buf.WriteString("\n")
		buf.WriteString("\treturn json.Marshal(res)")
		buf.WriteString("\n")
		buf.WriteString("}")
		buf.WriteString("\n")
	}

	if b.WeakUnmarshal {
		// We override the UnmarshalJSON method to be able to attempt to unmarshal scalar fields that should be
		// an array and to convert strings that should be put into numeric fields
		buf.WriteString("\n")
		buf.WriteString("// UnmarshalJSON implements the json.Unmarshaler interface and attempts weak decoding of fields")
		buf.WriteString("\n")
		buf.WriteString("func (b *Base) UnmarshalJSON(data []byte) error {")
		buf.WriteString("\n")

		for i, fieldName := range scalarFields {
			field := r.TopLevel[fieldName]
			if i == 0 {
				buf.WriteString(
					fmt.Sprintf(
						"\tresult := gjson.GetBytes(data, \"%s\")\n", field.Name,
					),
				)
			} else {
				buf.WriteString(
					fmt.Sprintf(
						"\tresult = gjson.GetBytes(data, \"%s\")\n", field.Name,
					),
				)
			}
			buf.WriteString("\tif result.Index > 0 {\n")
			if GoFieldType(field) == "[]string" || GoFieldType(field) == "int32" || GoFieldType(field) == "int64" {
				if GoFieldType(field) == "[]string" {
					buf.WriteString(
						fmt.Sprintf(
							"\t\tval_%s, err := decodeArray(\"%s\", result.Value())\n", field.Name, field.Path,
						),
					)
				} else if GoFieldType(field) == "int32" {
					buf.WriteString(
						fmt.Sprintf(
							"\t\tval_%s, err := decodeInt32(\"%s\", result.Value())\n", field.Name, field.Path,
						),
					)
				} else if GoFieldType(field) == "int64" {
					buf.WriteString(
						fmt.Sprintf(
							"\t\tval_%s, err := decodeInt64(\"%s\", result.Value())\n", field.Name, field.Path,
						),
					)
				}
				buf.WriteString("\t\tif err != nil {\n")
				buf.WriteString("\t\t\treturn err\n")
				buf.WriteString("\t\t}\n")
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = val_%s\n", field.FieldIdent().Pascal(), field.Name,
					),
				)
			} else if GoFieldType(field) == "string" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.String()\n", field.FieldIdent().Pascal(),
					),
				)
			} else if GoFieldType(field) == "time.Time" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.Time()\n", field.FieldIdent().Pascal(),
					),
				)
			} else if GoFieldType(field) == "float64" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.Float()\n", field.FieldIdent().Pascal(),
					),
				)
			} else if GoFieldType(field) == "bool" {
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = result.Bool()\n", field.FieldIdent().Pascal(),
					),
				)
			} else {
				buf.WriteString("\t\tvar raw []byte\n")
				buf.WriteString("\t\traw = data[result.Index:result.Index+len(result.Raw)]\n")
				buf.WriteString(fmt.Sprintf("\t\tvar aux %s\n", GoFieldType(field)))
				buf.WriteString("\t\terr := json.Unmarshal(raw, &aux)\n")
				buf.WriteString("\t\tif err != nil {\n")
				buf.WriteString("\t\t\treturn err\n")
				buf.WriteString("\t\t}\n")
				buf.WriteString(
					fmt.Sprintf(
						"\t\tb.%s = aux\n", field.FieldIdent().Pascal(),
					),
				)
			}
			buf.WriteString("\t}\n")
			buf.WriteString("\n")
		}

		// now we enumerate the object fields
		for i, fieldName := range objectFields {
			field := r.TopLevel[fieldName]
			if i == 0 {
				buf.WriteString(
					fmt.Sprintf(
						"\tresult_obj := gjson.GetBytes(data, \"%s\")\n", field.Name,
					),
				)
			} else {
				buf.WriteString(
					fmt.Sprintf(
						"\tresult_obj = gjson.GetBytes(data, \"%s\")\n", field.Name,
					),
				)
			}
			buf.WriteString("\tif result_obj.Index > 0 {\n")
			buf.WriteString("\t\tvar raw []byte\n")
			buf.WriteString("\t\traw = data[result_obj.Index:result_obj.Index+len(result_obj.Raw)]\n")
			buf.WriteString(fmt.Sprintf("\t\tvar aux %s\n", GoFieldType(field)))
			buf.WriteString("\t\terr := json.Unmarshal(raw, &aux)\n")
			buf.WriteString("\t\tif err != nil {\n")
			buf.WriteString("\t\t\treturn err\n")
			buf.WriteString("\t\t}\n")
			buf.WriteString(
				fmt.Sprintf(
					"\t\tb.%s = aux\n", field.FieldIdent().Pascal(),
				),
			)
			buf.WriteString("\t}\n")
		}

		buf.WriteString("\treturn nil\n")
		buf.WriteString("}")
		buf.WriteString("\n")

		// Create utility functions once used by weak unmarshal methods
		buf.WriteString("func decodeInt64(varName string, v interface{}) (int64, error) {")
		buf.WriteString("\n")
		buf.WriteString("\tdataVal := reflect.Indirect(reflect.ValueOf(v))\n")
		buf.WriteString("\tdataValKind := dataVal.Kind()\n")
		buf.WriteString("\tswitch dataValKind {\n")
		buf.WriteString("\tcase reflect.String:\n")
		buf.WriteString("\t\tif v == \"\" {\n")
		buf.WriteString("\t\t\treturn int64(0), nil\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t\tval, err := strconv.ParseInt(dataVal.String(), 10, 64)\n")
		buf.WriteString("\t\tif err != nil {\n")
		buf.WriteString("\t\t\treturn int64(0), fmt.Errorf(\"cannot parse var %s as int64\", varName)\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t\treturn int64(val), nil\n")
		buf.WriteString("\tcase reflect.Int, reflect.Int32, reflect.Int64:\n")
		buf.WriteString("\t\treturn int64(dataVal.Int()), nil\n")
		buf.WriteString("\tcase reflect.Float32, reflect.Float64:\n")
		buf.WriteString("\t\treturn int64(dataVal.Float()), nil\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tres, isOk := v.(int64)\n")
		buf.WriteString("\tif !isOk {\n")
		buf.WriteString("\t\treturn int64(0), fmt.Errorf(\"var %s is not of type int64\", varName)\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\treturn res, nil\n")
		buf.WriteString("}")
		buf.WriteString("\n")

		buf.WriteString("func decodeInt32(varName string, v interface{}) (int32, error) {")
		buf.WriteString("\n")
		buf.WriteString("\tdataVal := reflect.Indirect(reflect.ValueOf(v))\n")
		buf.WriteString("\tdataValKind := dataVal.Kind()\n")
		buf.WriteString("\tswitch dataValKind {\n")
		buf.WriteString("\tcase reflect.String:\n")
		buf.WriteString("\t\tif v == \"\" {\n")
		buf.WriteString("\t\t\treturn int32(0), nil\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t\tval, err := strconv.ParseInt(dataVal.String(), 10, 32)\n")
		buf.WriteString("\t\tif err != nil {\n")
		buf.WriteString("\t\t\treturn int32(0), fmt.Errorf(\"cannot parse var %s as int32\", varName)\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t\treturn int32(val), nil\n")
		buf.WriteString("\tcase reflect.Int, reflect.Int32, reflect.Int64:\n")
		buf.WriteString("\t\treturn int32(dataVal.Int()), nil\n")
		buf.WriteString("\tcase reflect.Float32, reflect.Float64:\n")
		buf.WriteString("\t\treturn int32(dataVal.Float()), nil\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tres, isOk := v.(int32)\n")
		buf.WriteString("\tif !isOk {\n")
		buf.WriteString("\t\treturn int32(0), fmt.Errorf(\"var %s is not of type int32\", varName)\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\treturn res, nil\n")
		buf.WriteString("}")
		buf.WriteString("\n")

		// Only arrays of strings exist in ECS
		buf.WriteString("func decodeArray(varName string, v interface{}) ([]string, error) {")
		buf.WriteString("\n")
		buf.WriteString("\tdataVal := reflect.Indirect(reflect.ValueOf(v))\n")
		buf.WriteString("\tdataValKind := dataVal.Kind()\n")
		buf.WriteString("\tswitch dataValKind {\n")
		buf.WriteString("\tcase reflect.String:\n")
		buf.WriteString("\t\treturn strings.Split(dataVal.String(), \",\"), nil\n")
		buf.WriteString("\tcase reflect.Int, reflect.Int64:\n")
		buf.WriteString("\t\tval := strconv.FormatInt(dataVal.Int(), 10)\n")
		buf.WriteString("\t\treturn []string{val}, nil\n")
		buf.WriteString("\tcase reflect.Uint, reflect.Uint64:\n")
		buf.WriteString("\t\tval := strconv.FormatUint(dataVal.Uint(), 10)\n")
		buf.WriteString("\t\treturn []string{val}, nil\n")
		buf.WriteString("\tcase reflect.Float32, reflect.Float64:\n")
		buf.WriteString("\t\tval := strconv.FormatFloat(dataVal.Float(), 'f', -1, 64)\n")
		buf.WriteString("\t\treturn []string{val}, nil\n")
		buf.WriteString("\tcase reflect.Array, reflect.Slice:\n")
		buf.WriteString("\t\tval := reflect.ValueOf(v)\n")
		buf.WriteString("\t\tsliceLen := val.Len()\n")
		buf.WriteString("\t\tresult := make([]string, 0, sliceLen)\n")
		buf.WriteString("\t\tfor i := 0; i < sliceLen; i++ {\n")
		buf.WriteString("\t\t\tresult = append(result, val.Index(i).Elem().String())\n")
		buf.WriteString("\t\t}\n")
		buf.WriteString("\t\treturn result, nil\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\tres, isOk := v.([]string)\n")
		buf.WriteString("\tif !isOk {\n")
		buf.WriteString("\t\treturn nil, fmt.Errorf(\"var %s is not of type string array\", varName)\n")
		buf.WriteString("\t}\n")
		buf.WriteString("\treturn res, nil\n")
		buf.WriteString("}")
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

// Execute implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (b *basic) Execute(root *ecsgen.Root) error {
	keys := []string{}

	// enumerate through for all implied objects
	// and sort them so the generation is deterministic
	for p, node := range root.Index {
		if node.IsObject() {
			keys = append(keys, p)
		}
	}

	sort.Strings(keys)

	// Create a buffer to write the source code to as we generate it
	// Using a bytes.Buffer over a strings.Builder because the go/parser
	// uses []byte in the parser.ParseFile function to parse sourcecode.
	buf := new(bytes.Buffer)

	// Add the generated comment and the package definition
	buf.WriteString("// Code generated by ecsgen; DO NOT EDIT.\n")
	buf.WriteString(fmt.Sprintf("package %s\n\n", b.PackageName))

	if b.WeakUnmarshal {
		// Add the gjson import since this is not handled by imports.Process, not sure why
		buf.WriteString("import \"github.com/tidwall/gjson\"\n")
	}

	// Add the top level Base type definition at the top of the file
	baseDef, err := b.CreateBase(root)
	if err != nil {
		return fmt.Errorf("error generating Base type definition: %v", err)
	}

	buf.WriteString(baseDef)
	buf.WriteString("\n")

	// Enumerate through all the objects, sorted by name alphabetically
	// and add their type definitions to the buffer
	for _, k := range keys {
		obj := root.Branch(k)
		code, err := b.ToGoCode(obj)
		if err != nil {
			return fmt.Errorf("error generating go code for %s: %v", k, err)
		}
		buf.WriteString(code)
	}

	// Create a new fileset and parse the generated Go code
	// this should catch any compile-time syntax errors we might have
	fs := token.NewFileSet()
	astFile, err := parser.ParseFile(fs, b.Filename, buf.Bytes(), parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing generated go code: %v", err)
	}

	// Format the Go code - this step is redundant, because the imports.Process
	// call below will also "pretty print", but I prefer to do it because
	// I'd rather have the formatted code have be validated by parser.ParseFile
	// first, before trying to format it.
	dstBuf := new(bytes.Buffer)
	err = format.Node(dstBuf, fs, astFile)
	if err != nil {
		return fmt.Errorf("error formatting generated go code: %v", err)
	}

	// // Now we will handle the imports
	imported, err := imports.Process(b.Filename, dstBuf.Bytes(), nil)
	if err != nil {
		return fmt.Errorf("error adding imports to generated go code: %v", err)
	}

	// Now write the resulting Go code to a file
	err = ioutil.WriteFile(filepath.Join(b.OutputDir, b.Filename), imported, 0644)
	if err != nil {
		return fmt.Errorf("error writing go code to file: %v", err)
	}

	return nil
}
