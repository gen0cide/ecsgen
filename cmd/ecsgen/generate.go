package main

import (
	"ecsgen/config"
	"ecsgen/loader"
	"fmt"
	"github.com/urfave/cli"
)

func init() {
	// create the loader config
	c, err := config.NewEmptyConfig()
	if err != nil {
		panic(err)
	}

	genConfig = c

	generateCommand = &cli.Command{
		Name:        "generate",
		Aliases:     []string{"g"},
		Usage:       "Use to translate ECS YAML definitions into a Go package.",
		Description: "Takes the input YAML definitions and translates into Idiomatic Go code.",
		Flags:       genConfig.CLIFlags(),
		Action:      generate,
	}
}

var (
	genConfig       *config.Config
	generateCommand *cli.Command
)

func generate(c *cli.Context) error {
	logger.Info("Running Generator")

	loader, err := loader.NewLoader(genConfig)
	if err != nil {
		return err
	}

	err = loader.Load()
	if err != nil {
		return err
	}

	root := loader.Root()

	generators, err := genConfig.Generators()
	if err != nil {
		return err
	}

	for _, g := range generators {
		err := g.Execute(root)
		if err != nil {
			return fmt.Errorf("error running %s generator: %v", g.ID(), err)
		}
	}

	return nil
}
