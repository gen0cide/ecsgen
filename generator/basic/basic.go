package basic

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
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
	"github.com/urfave/cli"
)

var (
	// ErrInvalidPackageName is thrown when a Go package name is either not specified or is not valid.
	ErrInvalidPackageName = errors.New("package name was either empty or an invalid go package identifier")

	// ErrInvalidOutputDir is thrown when the output directory does not exist.
	ErrInvalidOutputDir = errors.New("output directory was either blank or did not exist")
)

type basic struct {
	PackageName string
	OutputDir   string
	Pointers    bool
}

// New is a constructor for an empty debug output plugin.
func New() generator.Generator {
	return &basic{}
}

// ID implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (b *basic) ID() string {
	return "go_basic"
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
		&cli.BoolFlag{
			Name:        "use-pointers",
			Usage:       "Force the generator to use pointer types as best as possible.",
			EnvVars:     []string{"USE_POINTERS"},
			Destination: &b.Pointers,
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
		return "time.Duration"
	case n.Name == "args" && n.Definition.Type == "keyword":
		return "[]string"
	}

	// Find the right type!
	switch n.Definition.Type {
	case "keyword", "text", "ip", "geo_point":
		return "string"
	case "long":
		return "int64"
	case "integer":
		return "int32"
	case "float":
		return "float64"
	case "date":
		return "time.Time"
	case "boolean":
		return "bool"
	case "object":
		return "map[string]interface{}"
	default:
		panic(fmt.Errorf("no translation for %v (field %s)", n.Definition.Type, n.Name))
	}
}

// ToGoCode attempts to convert an ecsgen.Node into a Golang struct definition.
func ToGoCode(n *ecsgen.Node) (string, error) {
	if !n.IsObject() {
		return "", fmt.Errorf("node %s is not an object", n.Path)
	}

	scalarKeys := []string{}
	objectKeys := []string{}

	for key := range n.Children {
		scalarKeys = append(scalarKeys, key)
	}

	sort.Strings(scalarKeys)
	sort.Strings(objectKeys)

	buf := new(strings.Builder)

	buf.WriteString(fmt.Sprintf("// %s defines the object located at ECS path %s.", n.TypeIdent().Pascal(), n.Path))
	buf.WriteString("\n")
	buf.WriteString(fmt.Sprintf("type %s struct {", n.TypeIdent().Pascal()))
	buf.WriteString("\n")
	for _, k := range scalarKeys {
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
	buf.WriteString("}")
	buf.WriteString("\n")

	return buf.String(), nil
}

/*

TODO: Figure out how to make this work

1. Generate the list of base fields, each one gets their own file
2. Generate each file concurrently - write the output
3. Generate the base type - write the output

*/

type codegen struct {
	fs    *token.FileSet
	files map[string]*file
}

type file struct {
	root    *ecsgen.Node
	codeAST *ast.File
	data    *bytes.Buffer
}

func (b *basic) Generate(n *ecsgen.Node, finchan chan struct{}, errchan chan error) {
	buf := new(bytes.Buffer)
	_ = buf
}

// Execute implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (b *basic) Execute(root *ecsgen.Root) error {
	keys := []string{}

	// enumerate through for all implied objects
	for p, node := range root.Index {
		if node.IsObject() {
			keys = append(keys, p)
		}
	}

	sort.Strings(keys)

	buf := new(bytes.Buffer)

	buf.WriteString("// Code generated by ecsgen; DO NOT EDIT.\n")
	buf.WriteString(fmt.Sprintf("package %s\n\n", b.PackageName))

	for _, k := range keys {
		obj := root.Branch(k)
		code, err := ToGoCode(obj)
		if err != nil {
			return fmt.Errorf("error generating go code for %s: %v", k, err)
		}
		buf.WriteString(code)
	}

	fs := token.NewFileSet()
	astFile, err := parser.ParseFile(fs, "generated_definitions.go", buf.Bytes(), parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing generated go code: %v", err)
	}

	dstBuf := new(bytes.Buffer)
	err = format.Node(dstBuf, fs, astFile)
	if err != nil {
		return fmt.Errorf("error formatting generated go code: %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(b.OutputDir, "generated_definitions.go"), dstBuf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing output definitions: %v", err)
	}

	return nil
}
