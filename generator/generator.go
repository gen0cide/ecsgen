package generator

import (
	"fmt"
	"strings"

	"github.com/gen0cide/ecsgen"
	"github.com/urfave/cli"
)

// The Generator interface is used to define a way to create multiple output
// plugins that allow ecsgen to create various code off of a single loaded schema.
type Generator interface {
	// this ID must be a snake case string and must be unique.
	ID() string

	// CLIFlags is called by the executable in order to
	// inject the flags into the CLI framework. All the CLI flags
	// will have their Name and EnvVars values modified to be prefixed.
	// For example, a Generator who's ID() value is "basic" and has
	// a cli.StringFlag with a name of "package-name" and EnvVars of "PACKAGE_NAME"
	// will have the resulting cli access:
	// --opt-basic-package-name (cli flag)
	// ECSGEN_OPT_BASIC_PACKAGE_NAME (env var)
	CLIFlags() []cli.Flag

	// Validate ensures that configuration is valid
	// for the the given code generator.
	Validate() error

	// Execute will be called by the application when the loader has
	// successfully loaded the package.
	Execute(r *ecsgen.Root) error
}

// shimCLIFlags is used to shim each generators CLI flags to prefix them with the correct
// flag names, as well as environment variable prefixes.
func shimCLIFlags(g Generator) []cli.Flag {
	originalFlags := g.CLIFlags()
	if len(originalFlags) == 0 {
		return []cli.Flag{}
	}

	newFlags := []cli.Flag{}
	pluginID := ecsgen.NewIdentifier(g.ID())

	for _, flag := range originalFlags {
		var newFlag cli.Flag

		// figure out the type of the flag and shim it
		switch tflag := flag.(type) {
		case *cli.BoolFlag:
			newFlag = &cli.BoolFlag{
				Name:        prefixName(pluginID, tflag.Name),
				Usage:       tflag.Usage,
				EnvVars:     prefixEnvVars(pluginID, tflag.EnvVars),
				Destination: tflag.Destination,
			}
		case *cli.StringFlag:
			newFlag = &cli.StringFlag{
				Name:        prefixName(pluginID, tflag.Name),
				Usage:       tflag.Usage,
				EnvVars:     prefixEnvVars(pluginID, tflag.EnvVars),
				Destination: tflag.Destination,
				Value:       tflag.Value,
			}
		case *cli.StringSliceFlag:
			newFlag = &cli.StringSliceFlag{
				Name:        prefixName(pluginID, tflag.Name),
				Usage:       tflag.Usage,
				EnvVars:     prefixEnvVars(pluginID, tflag.EnvVars),
				Destination: tflag.Destination,
				Value:       tflag.Value,
			}
		default:
			panic(fmt.Errorf("output generator does not yet implement flags of type %T", tflag))
		}
		newFlags = append(newFlags, newFlag)
	}

	return newFlags
}

func prefixEnvVars(id ecsgen.Identifier, vars []string) []string {
	// short circuit if vars is 0 length
	if len(vars) == 0 {
		return vars
	}

	newvars := make([]string, len(vars))

	for idx, val := range vars {
		newvars[idx] = strings.Join([]string{"ECSGEN_OPT", id.Screaming(), val}, "_")
	}

	return newvars
}

func prefixName(id ecsgen.Identifier, name string) string {
	return strings.Join([]string{"opt", id.Command(), name}, "-")
}
