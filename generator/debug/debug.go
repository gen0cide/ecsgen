package debug

import (
	"github.com/gen0cide/ecsgen"
	"github.com/gen0cide/ecsgen/generator"
	"github.com/urfave/cli"
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
	return nil
}
