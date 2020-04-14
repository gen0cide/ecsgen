package main

import (
	"os"

	"ecsgen"
	"github.com/gen0cide/genolog"
	"github.com/gen0cide/genolog/pretty"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	// globals
	debugLog = false
	logger   genolog.Logger
)

func main() {
	logger = pretty.NewPrettyLogger("ECS", "gen", nil)

	app := cli.NewApp()
	app.Name = "ecsgen"
	app.Usage = "Generate Go types based off Elastic Common Schema (ECS) definitions."
	app.UsageText = "ecsgen [--version|-v] [--debug|-d] [--help|-h] COMMAND [COMMAND_OPTIONS]"
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "Alex Levinson",
			Email: "gen0cide.threats@gmail.com",
		},
	}
	app.Version = ecsgen.Version
	app.Description = "More information on this tool can be found at https://github.com/gen0cide/ecsgen."
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Aliases:     []string{"d"},
			Usage:       "Enables more verbose, debug logging.",
			EnvVars:     []string{"ECSGEN_DEBUG"},
			Value:       false,
			Destination: &debugLog,
		},
	}
	app.Before = func(c *cli.Context) error {
		if debugLog {
			logger.LogrusLogger().SetLevel(logrus.DebugLevel)
		}

		return nil
	}
	app.Commands = []*cli.Command{
		generateCommand,
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatalw("exiting with error", "error", err)
	}
}
