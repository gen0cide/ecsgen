package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"ecsgen/generator"
	"ecsgen/generator/avro"
	"ecsgen/generator/basic"
	"ecsgen/generator/debug"
	"github.com/urfave/cli"
)

var (
	// ErrInvalidSourceFile is thrown when the source file cannot be resolved
	ErrInvalidSourceFile = errors.New("source file must be specified")

	// ErrNoDefinitionsInSourceFile is thrown when the source directory does not contain any valid
	// ECS definitions.
	ErrNoDefinitionsInSourceFile = errors.New("source directory does not contain any valid ecs definitions")
)

var (
	// list of the builtin generators
	builtinGenerators = []generator.Generator{
		debug.New(),
		basic.New(),
		avro.New(),
	}
)

// Config holds the parameters needed for proper generation of Go code.
type Config struct {
	SourceFile string

	whitelist  *cli.StringSlice
	blacklist  *cli.StringSlice
	generators *cli.StringSlice
	registry   generator.Registry
}

// NewEmptyConfig is a constructor for an empty Config object.
func NewEmptyConfig() (*Config, error) {
	registry := generator.NewRegistry()

	for _, x := range builtinGenerators {
		err := registry.Register(x)
		if err != nil {
			return nil, fmt.Errorf("could not register output plugin: %v", err)
		}
	}

	return &Config{
		whitelist:  cli.NewStringSlice(),
		blacklist:  cli.NewStringSlice(),
		generators: cli.NewStringSlice(),
		registry:   registry,
	}, nil
}

// CLIFlags is a helper to automatically set fields within a Config object
// based on CLI flags using the github.com/urfave/cli framework.
func (c *Config) CLIFlags() []cli.Flag {
	pluginNames := []string{}
	for _, x := range c.registry.All() {
		pluginNames = append(pluginNames, x.ID())
	}

	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "source-file",
			Usage:       "Path to the generated ecs_flat.yml file containing ECS definitions.",
			EnvVars:     []string{"ECSGEN_SOURCE_FILE"},
			Required:    true,
			Destination: &c.SourceFile,
		},
		&cli.StringSliceFlag{
			Name:        "whitelist",
			Usage:       "Regular expression that denotes which ECS keys to allow into the model. (Can be used multiple times).",
			EnvVars:     []string{"ECSGEN_WHITELIST_VALUE"},
			Value:       c.whitelist,
			Destination: c.whitelist,
		},
		&cli.StringSliceFlag{
			Name:        "blacklist",
			Usage:       "Regular expression that denotes which ECS keys to explicitly forbid into the model. (Can be used multiple times).",
			EnvVars:     []string{"ECSGEN_BLACKLIST_VALUE"},
			Value:       c.blacklist,
			Destination: c.blacklist,
		},
		&cli.StringSliceFlag{
			Name:        "output-plugin",
			Usage:       fmt.Sprintf("Enable an output generator plugin. Can be used multiple times. Possible values: %s", strings.Join(pluginNames, ", ")),
			EnvVars:     []string{"ECSGEN_OUTPUT_PLUGIN"},
			Value:       c.generators,
			Destination: c.generators,
		},
	}

	flags = append(flags, c.registry.CLIFlags()...)
	return flags
}

// Validate validates that the configuration has expected parameters.
func (c *Config) Validate() error {
	// Check the Source Directory
	// is it assigned?
	if c.SourceFile == "" {
		return ErrInvalidSourceFile
	}

	// Is it a valid path?
	dir, err := os.Stat(c.SourceFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("could not locate the source file: %v", err)
		}

		return fmt.Errorf("error locating specified source file: %v", err)
	}

	// is it a valid file?
	if dir.IsDir() {
		return fmt.Errorf("specified source file path was a directory, not a file")
	}

	// check to make sure the whitelist is valid
	if _, err := c.Whitelist(); err != nil {
		return fmt.Errorf("error parsing whitelist parameter: %v", err)
	}

	// check to make sure the whitelist is valid
	if _, err := c.Blacklist(); err != nil {
		return fmt.Errorf("error parsing blacklist parameter: %v", err)
	}

	// verify output plugins
	if len(c.generators.Value()) == 0 {
		return fmt.Errorf("did not specify any output generators")
	}

	// ensure any specified output plugins actually exist
	for _, x := range c.generators.Value() {
		generator, err := c.registry.Get(x)
		if err != nil {
			pluginNames := []string{}
			for _, x := range c.registry.All() {
				pluginNames = append(pluginNames, x.ID())
			}
			return fmt.Errorf("%s is not a valid plugin name. valid options: %s", x, strings.Join(pluginNames, ", "))
		}

		// check that the configuration is valid for any enabled plugin
		err = generator.Validate()
		if err != nil {
			return fmt.Errorf("error in output plugin %s: %v", generator.ID(), err)
		}
	}

	return nil
}

// Generators returns the set of enabled generators for the config.
func (c *Config) Generators() ([]generator.Generator, error) {
	ret := []generator.Generator{}

	for _, x := range c.generators.Value() {
		generator, err := c.registry.Get(x)
		if err != nil {
			pluginNames := []string{}
			for _, x := range c.registry.All() {
				pluginNames = append(pluginNames, x.ID())
			}
			return ret, fmt.Errorf("%s plugin could not be loaded. valid options: %s", x, strings.Join(pluginNames, ", "))
		}
		ret = append(ret, generator)
	}

	return ret, nil
}
