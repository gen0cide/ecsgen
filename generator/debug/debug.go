package debug

import (
	"fmt"
	"strings"

	"github.com/gen0cide/ecsgen"
	"github.com/gen0cide/ecsgen/generator"
	"github.com/urfave/cli/v2"
)

type debug struct {
}

// New is a constructor for an empty debug output plugin.
func New() generator.Generator {
	return &debug{}
}

// ID implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (d *debug) ID() string {
	return "debug"
}

// CLIFlags implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (d *debug) CLIFlags() []cli.Flag {
	return []cli.Flag{}
}

// Validate implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (d *debug) Validate() error {
	return nil
}

// Execute implements the generator.Generator interface.
// Package: github.com/gen0cide/ecsgen/generator
func (d *debug) Execute(r *ecsgen.Root) error {
	walkFn := func(n *ecsgen.Node) error {
		level := len(strings.Split(n.Path, ".")) - 1
		indent := new(strings.Builder)
		for i := 0; i < level; i++ {
			indent.WriteString("\t")
		}
		if n.IsObject() {
			fmt.Printf("%s[OBJECT] %s\n", indent.String(), n.Path)
			return nil
		}

		fmt.Printf("%s(field) %s\n", indent.String(), n.Path)
		return nil
	}

	err := ecsgen.Walk(r, walkFn)
	if err != nil {
		return fmt.Errorf("error walking tree: %v", err)
	}

	return nil
}
