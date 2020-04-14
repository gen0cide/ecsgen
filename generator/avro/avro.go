package avro

import (
	"bytes"
	"errors"
	"fmt"
	"encoding/json"

	"os"
	"regexp"
	"io/ioutil"

	"ecsgen"
	"ecsgen/generator"
	"github.com/urfave/cli"
)



var (
	// ErrInvalidPackageName is thrown when a Go package name is either not specified or is not valid.
	ErrInvalidPackageName = errors.New("package name was either empty or an invalid go package identifier")

	// ErrInvalidOutputDir is thrown when the output directory does not exist.
	ErrInvalidOutputDir = errors.New("output directory was either blank or did not exist")

	// ErrInvalidOwner is thrown when the owner is not specified.
	ErrInvalidOwner = errors.New("owner was blank")

	// ErrInvalidNamespace is thrown when the namespace is not specified.
	ErrInvalidNamespace = errors.New("namespace was blank")

	// ErrInvalidName is thrown when the name is not specified.
	ErrInvalidName = errors.New("name was blank")
)

type avro struct {
	PackageName string
	OutputDir   string
	Namespace string
	Name string
	Owner string
	Pointers    bool
}

// New is a constructor for an empty debug output plugin.
func New() generator.Generator {
	return &avro{}
}

// ID implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (a *avro) ID() string {
	return "go_avro"
}

// CLIFlags implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (a *avro) CLIFlags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "package-name",
			Usage:       "Name of the Go package for the generated code.",
			EnvVars:     []string{"PACKAGE_NAME"},
			Destination: &a.PackageName,
		},
		&cli.StringFlag{
			Name:        "output-dir",
			Usage:       "Path to the directory where the generated code should be written.",
			EnvVars:     []string{"OUTPUT_DIR"},
			Destination: &a.OutputDir,
		},
		&cli.StringFlag{
			Name:        "owner",
			Usage:       "The individual responsible for this schema",
			Destination: &a.Owner,
		},
		&cli.StringFlag{
			Name:        "namespace",
			Usage:       "This identifies the namespace in which the object lives. Essentially, this is meant to be a URI that has meaning to you and your organization. It is used to differentiate one schema type from another should they share the same name.",
			Destination: &a.Namespace,
		},
		&cli.StringFlag{
			Name:        "name",
			Usage:       "This is the schema name which, when combined with the namespace, uniquely identifies the schema within the store. In the above example, the fully qualified name for the schema is com.example.FullName.",
			Destination: &a.Name,
		},
	}
}

// Validate implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (a *avro) Validate() error {

	if a.Owner == "" {
		return ErrInvalidOwner
	}
	if a.Namespace == "" {
		return ErrInvalidNamespace
	}
	if a.Name == "" {
		return ErrInvalidName
	}


	// Check the Output Directory
	// is it assigned?
	if a.OutputDir == "" {
		return ErrInvalidOutputDir
	}

	// Is it a valid path?
	dir, err := os.Stat(a.OutputDir)
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
	if !pkgRegex.MatchString(a.PackageName) {
		return ErrInvalidPackageName
	}

	return nil
}

// GoFieldType returns the Go type to be used in the Go struct field type definition.
func GoFieldType(n *ecsgen.Node) string {
	// create a buffer to determine type
	//typeBuf := new(bytes.Buffer)

	//if n.IsArray() {
	//	typeBuf.WriteString("[]")
	//}

	// if Node is an Object, we need to return this object's type. For example,
	// Node("client.nat") needs to return "ClientNAT" as it's Go type.
	//if n.IsObject() {
	//	typeBuf.WriteString(n.TypeIdent().Pascal())
	//	return typeBuf.String()
	//}

	// Special cases denoted by the ECS developers.
	switch {
	case n.Name == "duration" && n.Definition.Type == "long":
		return "long"
	case n.Name == "args" && n.Definition.Type == "keyword":
		return "string"
	}

	// Check date type, review avro types and correspond (make note)
	//@TODO
	//Check if yaml has description -> (watch tower has a short descirption -> avro does not)
	// Find the right type!
	switch n.Definition.Type {
	case "keyword", "text", "ip", "geo_point":
		return "string"
	case "long":
		return "long"
	case "integer":
		return "int"
	case "float":
		return "double"
	case "date":
		return "string"
	case "boolean":
		return "boolean"
	//	Will need to handle this differently
	case "object":
		return "string"
	default:
		panic(fmt.Errorf("no translation for %v (field %s)", n.Definition.Type, n.Name))
	}
}

// ToGoCode attempts to convert an ecsgen.Node into a Golang struct definition.
//{
//"date_created": "2018-07-16 23:20:27 -0753",
//"fields": [
//{
//"default": null,
//"name": "hello",
//"type": [
//"null",
//"long"
//]
//},
//{
//"default": null,
//"name": "simple",
//"type": [
//"null",
//"string"
//]
//}
//],
//"name": "event_04544554",
//"namespace": "082402.kahwee",
//"owner": "",
//"schemaVersion": 2,
//"schema_id": 1,
//"type": "record"
//}
//{
//"type": "record",
//"namespace": "com.example",
//"name": "FullName",
//"fields": [
//{ "name": "first", "type": "string" },
//{ "name": "last", "type": "string" }
//]
//}

//{
//"type": "record",
//"namespace": "foz",
//"name": "trip_events",
//"fields": [
//{"name": "id"        , "type": "long"              , "default": 0   },
//{"name": "client_id" , "type": ["null"   , "long"] , "default": null},
//{"name": "driver_id" , "type": ["null"   , "long"] , "default": null},
//{"name": "event"     , "type": ["string" , "null"] , "default": ""  },
//{"name": "lat"       , "type": ["null"   , "float"], "default": null},
//{"name": "lng"       , "type": ["null"   , "float"], "default": null}
//],
//"owner": "Jane Engineer"
//}

type Schema struct {
	Type string `json:"type"`
	Namespace          string `json:"namespace"`
	Name        string `json:"name"`
	Fields        []*Field `json:"fields" yaml:"type,omitempty" ecs:"agent.type"`
	Owner     string `json:"owner"`
}

// must validate the NestedTypes and Type only one can exist
type Field struct {
	Name string `json:"name"`
	Parent string //we need to keep track of the name of the parent of the nested types
	NestedTypes          []*Field `json:"fields,omitempty"`
	Type string `json:"type"`
	Defined bool //defaults to false
	//Doc string `json:"doc"`
	//Default        []string `json:"default"` // this might not be needed..
}

func (f *Field) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})
	ret["name"] = f.Name
	//ret["doc"] = f.Doc
	var typeArray []interface{}
	typeArray = append(typeArray, "null")

	if f.Defined {
		ret["type"] = append(typeArray, f.Name)
	} else {
		switch f.Type {
		case "record":
			fieldType := map[string]interface{}{
				"fields": f.NestedTypes,
				"type":   "record",
				"name":   f.Name,
			}
			ret["type"] = append(typeArray, fieldType)
		default:
			//ret["type"] = "string"
			ret["type"] = f.Type
		}
	}
	return json.Marshal(ret)
}

func customPostOrder(root *ecsgen.Node, parent string, fieldsDefined map[string]bool) (*Field){
	field := Field{
		Name: root.Name,
	}

	if root.Definition != nil{
		//field.Doc = root.Definition.Description

		//field.Type = root.Definition.Type
		field.Type = GoFieldType(root)
	} else {
		//check to see if we have defined this record previously
		exists := fieldsDefined[root.Name]
		fmt.Println(exists)
		if exists {
			field.Defined = true
		}else {
			field.Parent = parent
			field.Type = "record"
			fields := []*Field{}
			for _, v := range root.Children {
				nestedField := customPostOrder(v, root.Name, fieldsDefined)
				fields = append(fields, nestedField)

			}
			field.NestedTypes = fields
			fieldsDefined[root.Name] = true
			//fmt.Println(fieldsDefined[root.Name])
		}
	}
	return &field
}

func (a *avro) Generate(n *ecsgen.Node, finchan chan struct{}, errchan chan error) {
	buf := new(bytes.Buffer)
	_ = buf
}

// Execute implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (a *avro) Execute(root *ecsgen.Root) error {
	schema := Schema{
		Type: "record",
		Name: a.Name,
		Namespace: a.Namespace,
		Owner: a.Owner,
	}
	//fmt.Println(schema)

	fields := []*Field{}
	fieldsDefined := make(map[string]bool)
	for _, node := range root.TopLevel{
		if node.IsObject() {
			nestedField := customPostOrder(node, "", fieldsDefined)
			fields = append(fields, nestedField)
		}
	}
	schema.Fields = fields

	file, err := os.OpenFile("output.json", os.O_CREATE, 0644)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	//jsonData, err := json.Marshal(schema)
	jsonData, err := json.MarshalIndent(schema, "", "  ")
	fmt.Println(string(jsonData), err)
	if err != nil {
		fmt.Println(err)
		return err
	}
	//fmt.Println(jsonData)
	ioutil.WriteFile("avro_schema_2.json", jsonData, os.ModePerm)

	return nil
}
